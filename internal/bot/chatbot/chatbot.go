package chatbot

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
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

func getMuffinResponse(ctx context.Context, s *discordgo.Session, question string) (string, error) {
	var data []repository.Text
	var result string
	x := rand.Intn(10)

	cur, err := repository.GetDatabase().Texts.Find(context.TODO(), bson.D{{Key: "persona", Value: "muffin"}})
	if err != nil {
		return "살려주ㅅ세요", err
	}

	defer cur.Close(context.TODO())

	if err = cur.All(context.TODO(), &data); err != nil {
		return "살려주ㅅ세요", err
	}

	learnData, err := repository.GetDatabase().Knowledge.GetByCommand(ctx, question)
	if err != nil {
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

func getAIResponse(ctx context.Context, c *Chatbot, user *discordgo.User, question string, attachments ...*discordgo.MessageAttachment) (string, error) {
	const twelveHours = 43_200

	dbUser, err := repository.GetDatabase().Users.Get(ctx, user.ID)
	if err != nil {
		return "살려주ㅅ세요", err
	}

	if err := repository.GetDatabase().Chats.FindOne(context.TODO(), repository.Chat{UserID: user.ID}).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			if _, err = repository.GetDatabase().Chats.Create(ctx, user.ID, fmt.Sprintf("새로운 채팅 %06d", rand.Intn(999999))); err != nil {
				return "살려주ㅅ세요", err
			}
		} else {
			return "살려주ㅅ세요", err
		}
	}

	if dbUser.CreateNewChatAfter12Hours {
		timestamp, err := repository.GetDatabase().Memory.GetLastMemoryTimestamp(ctx, dbUser.ChatID)
		if err != nil {
			return "살려주세요", err
		}

		if time.Now().Unix()-timestamp > twelveHours {
			result, err := repository.GetDatabase().Chats.Create(ctx, user.ID, fmt.Sprintf("새로운 채팅 %06d", rand.Intn(999999)))
			if err != nil {
				return "살려주ㅅ세요", err
			}

			dbUser.ChatID = result.InsertedID.(bson.ObjectID)
		}
	}

	chat, err := c.GetChat(ctx, user, dbUser.ChatID)
	if err != nil {
		return "살려주ㅅ세요", err
	}

	var parts []genai.Part
	var files []repository.File

	if len(attachments) != 0 {
		genaiFiles, err := getFiles(c.Gemini, &attachments)
		if err != nil {
			return "살려주ㅅ세요", err
		}

		for _, file := range genaiFiles {
			parts = append(parts, *genai.NewPartFromFile(*file))
			files = append(files, repository.File{URI: file.URI, MIMEType: file.MIMEType})
		}
	}

	parts = append(parts, *genai.NewPartFromText(question))

	result, err := chat.SendMessage(context.TODO(), parts...)
	if err != nil {
		return "살려주ㅅ세요", err
	}

	resultText := result.Text()
	if err = repository.GetDatabase().Memory.Save(ctx, dbUser.ChatID, user.ID, question, resultText, files); err != nil {
		return "살려주ㅅ세요", err
	}

	log.Printf("%s TOKEN: %d", user.ID, result.UsageMetadata.PromptTokenCount)

	return fmt.Sprintf("%s\n`해당 문장은 AI가 생성한 것이에요.`", resultText), nil
}

func (c *Chatbot) GetResponse(ctx context.Context, user *discordgo.User, question string, attachments ...*discordgo.MessageAttachment) (string, error) {
	mode, err := repository.GetDatabase().Users.GetUserChattingMode(ctx, user.ID)
	if err != nil {
		return "살려주ㅅ세요", err
	}

	switch mode {
	case repository.ChattingMuffinMode:
		return getMuffinResponse(ctx, c.s, question)
	default:
		return getAIResponse(ctx, c, user, question, attachments...)
	}
}
