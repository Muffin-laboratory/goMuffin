package chatbot

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/genai"
)

type Chatbot struct {
	Gemini       *genai.Client
	systemPrompt string
	s            *discordgo.Session
}

var instance *Chatbot

func Make(s *discordgo.Session) error {
	gemini, err := genai.NewClient(context.TODO(), &genai.ClientConfig{
		APIKey:  configs.GetConfig().Chatbot.Gemini.Token,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return err
	}

	instance = &Chatbot{
		Gemini: gemini,
		s:      s,
	}

	prompt, err := loadPrompt()
	if err != nil {
		return err
	}

	instance.systemPrompt = prompt
	return nil
}

func GetChatBot() *Chatbot {
	return instance
}

func (c *Chatbot) ReloadPrompt() error {
	prompt, err := loadPrompt()
	if err != nil {
		return err
	}

	c.systemPrompt = prompt
	return nil
}

func (c *Chatbot) GetPrompt() string {
	return c.systemPrompt
}

func getMuffinResponse(s *discordgo.Session, question string) (string, error) {
	var learnData []databases.Knowledge
	var data []databases.Text
	var result string
	x := rand.Intn(10)

	muffinCur, err := databases.GetDatabase().Texts.Find(context.TODO(), bson.D{{Key: "persona", Value: "muffin"}})
	if err != nil {
		return "살려주ㅅ세요", err
	}
	learnCur, err := databases.GetDatabase().Knowledge.Find(context.TODO(), bson.D{{Key: "command", Value: question}})
	if err != nil {
		return "살려주ㅅ세요", err
	}

	defer muffinCur.Close(context.TODO())
	defer learnCur.Close(context.TODO())

	if err = muffinCur.All(context.TODO(), &data); err != nil {
		return "살려주ㅅ세요", err
	}

	if err = learnCur.All(context.TODO(), &learnData); err != nil {
		return "살려주ㅅ세요", err
	}

	if x > 2 && len(learnData) != 0 {
		data := learnData[rand.Intn(len(learnData))]
		user, _ := s.User(data.UserID)

		result =
			fmt.Sprintf("%s\n%s", data.Result, utils.InlineCode(fmt.Sprintf("%s님이 알려주셨어요.", user.Username)))
	} else {
		result = data[rand.Intn(len(data))].Text
	}
	return result, nil
}

func getAIResponse(c *Chatbot, user *discordgo.User, question string) (string, error) {
	var dbUser databases.User

	if err := databases.GetDatabase().Users.FindOne(context.TODO(), databases.User{
		UserID: user.ID,
	}).Decode(&dbUser); err != nil {
		return "살려주ㅅ세요", err
	}

	if err := databases.GetDatabase().Chats.FindOne(context.TODO(), databases.Chat{UserId: user.ID}).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = databases.GetDatabase().Chats.CreateChat(user.ID, "새로운 채팅")
			fmt.Println(err)
			if err != nil {
				return "살려주ㅅ세요", err
			}
		} else {
			return "살려주ㅅ세요", err
		}
	}

	contents, err := databases.GetDatabase().Memory.Get(dbUser.ChatID)
	if err != nil {
		return "AI에 문제가 생겼ㅇ어요.", err
	}

	prompt, err := makePrompt(c.systemPrompt, user)
	if err != nil {
		return "살려주ㅅ세요", err
	}

	contents = append(contents, genai.NewContentFromText(question, genai.RoleUser))
	result, err := c.Gemini.Models.GenerateContent(context.TODO(), configs.GetConfig().Chatbot.Gemini.Model, contents, &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(prompt, genai.RoleUser),
	})
	if err != nil {
		return "AI에 문제가 생겼ㅇ어요.", err
	}

	resultText := result.Text()
	if err = databases.GetDatabase().Memory.Save(&databases.Memory{
		UserID:  user.ID,
		Content: question,
		Answer:  resultText,
		ChatID:  dbUser.ChatID,
	}); err != nil {
		return "살려주ㅅ세요", err
	}

	log.Printf("%s TOKEN: %d", user.ID, result.UsageMetadata.PromptTokenCount)

	return resultText, nil
}

func (c *Chatbot) GetResponse(user *discordgo.User, question string) (string, error) {
	mode, err := databases.GetDatabase().Users.GetUserChattingMode(user.ID)
	if err != nil {
		return "살려주ㅅ세요", err
	}

	switch mode {
	case databases.ChattingMuffinMode:
		return getMuffinResponse(c.s, question)
	default:
		return getAIResponse(c, user, question)
	}
}
