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
	"google.golang.org/genai"
)

type Chatbot struct {
	Mode         ChatbotMode
	Gemini       *genai.Client
	systemPrompt string
	s            *discordgo.Session
}

var ChatBot *Chatbot

func New(s *discordgo.Session) error {
	gemini, err := genai.NewClient(context.TODO(), &genai.ClientConfig{
		APIKey:  configs.Config.Chatbot.Gemini.Token,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return err
	}

	ChatBot = &Chatbot{
		Mode:   ChatbotAI,
		Gemini: gemini,
		s:      s,
	}

	prompt, err := loadPrompt()
	if err != nil {
		return err
	}

	ChatBot.systemPrompt = prompt
	return nil
}

func (c *Chatbot) SetMode(mode ChatbotMode) *Chatbot {
	c.Mode = mode
	return c
}

func (c *Chatbot) SwitchMode() *Chatbot {
	switch c.Mode {
	case ChatbotAI:
		c.SetMode(ChatbotMuffin)
	case ChatbotMuffin:
		c.SetMode(ChatbotAI)
	}
	return c
}

func (c *Chatbot) ModeString() string {
	switch c.Mode {
	case ChatbotAI:
		return "AI모드"
	case ChatbotMuffin:
		return "머핀 모드"
	default:
		return "알 수 없음"
	}
}

func (c *Chatbot) ReloadPrompt() error {
	prompt, err := loadPrompt()
	if err != nil {
		return err
	}

	c.systemPrompt = prompt
	return nil
}

func getMuffinResponse(s *discordgo.Session, question string) (string, error) {
	var data []databases.Text
	var learnData []databases.Learn
	var result string
	x := rand.Intn(10)

	muffinCur, err := databases.Database.Texts.Find(context.TODO(), bson.D{{Key: "persona", Value: "muffin"}})
	if err != nil {
		return "살려주ㅅ세요", err
	}
	learnCur, err := databases.Database.Learns.Find(context.TODO(), bson.D{{Key: "command", Value: question}})
	if err != nil {
		return "살려주ㅅ세요", err
	}

	defer muffinCur.Close(context.TODO())
	defer learnCur.Close(context.TODO())

	muffinCur.All(context.TODO(), &data)
	learnCur.All(context.TODO(), &learnData)

	if x > 2 && len(learnData) != 0 {
		data := learnData[rand.Intn(len(learnData))]
		user, _ := s.User(data.UserId)

		result =
			fmt.Sprintf("%s\n%s", data.Result, utils.InlineCode(fmt.Sprintf("%s님이 알려주셨어요.", user.Username)))
	} else {
		result = data[rand.Intn(len(data))].Text
	}
	return result, nil
}

func getAIResponse(s *discordgo.Session, c *Chatbot, user *discordgo.User, question string) (string, error) {
	var data []databases.Learn

	x := rand.Intn(10)

	cur, err := databases.Database.Learns.Find(context.TODO(), bson.D{{Key: "command", Value: question}})
	if err != nil {
		return "살려주ㅅ세요", err
	}

	defer cur.Close(context.TODO())
	cur.All(context.TODO(), &data)

	if x == 10 && len(data) != 0 {
		data := data[rand.Intn(len(data))]
		user, _ := s.User(data.UserId)
		return fmt.Sprintf("%s\n%s", data.Result, utils.InlineCode(fmt.Sprintf("%s님이 알려주셨어요.", user.Username))), nil
	}

	contents, err := GetMemory(user.ID)
	if err != nil {
		ChatBot.Mode = ChatbotMuffin
		return "AI에 문제가 생겼ㅇ어요.", err
	}

	contents = append(contents, genai.NewContentFromText(question, genai.RoleUser))
	result, err := ChatBot.Gemini.Models.GenerateContent(context.TODO(), configs.Config.Chatbot.Gemini.Model, contents, &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(makePrompt(c.systemPrompt, user), genai.RoleUser),
	})
	if err != nil {
		ChatBot.Mode = ChatbotMuffin
		return "AI에 문제가 생겼ㅇ어요.", err
	}

	resultText := result.Text()
	err = SaveMemory(&databases.Memory{
		UserId:  user.ID,
		Content: question,
		Answer:  resultText,
	})
	if err != nil {
		return "AI에 문제가 생겼ㅇ어요.", err
	}

	log.Printf("%s TOKEN: %d", user.ID, result.UsageMetadata.PromptTokenCount)

	return resultText, nil
}

func (c *Chatbot) GetResponse(user *discordgo.User, question string) (string, error) {
	switch c.Mode {
	case ChatbotMuffin:
		return getMuffinResponse(c.s, question)
	default:
		return getAIResponse(c.s, c, user, question)
	}
}
