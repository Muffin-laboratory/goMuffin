package chatbot

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/genai"
)

type Chatbot struct {
	Mode   ChatbotMode
	config *genai.GenerateContentConfig
	Gemini *genai.Client
	s      *discordgo.Session
}

var ChatBot *Chatbot

func New(s *discordgo.Session) {
	gemini, err := genai.NewClient(context.TODO(), &genai.ClientConfig{
		APIKey:  configs.Config.Chatbot.Gemini.Token,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalln(err)
	}

	ChatBot = &Chatbot{
		Mode:   ChatbotDefault,
		Gemini: gemini,
		s:      s,
	}

	bin, err := os.ReadFile(configs.Config.Chatbot.Gemini.PromptPath)
	if err != nil {
		log.Fatalln(err)
	}

	ChatBot.config = &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(string(bin), genai.RoleUser),
	}
}

func (c *Chatbot) SetMode(mode ChatbotMode) *Chatbot {
	c.Mode = mode
	return c
}

func getDefaultResponse(s *discordgo.Session, question string) string {
	var data []databases.Text
	var learnData []databases.Learn
	var result string

	ch := make(chan int)
	x := rand.Intn(10)

	// 머핀 데이터
	go func() {
		cur, err := databases.Database.Texts.Find(context.TODO(), bson.D{{Key: "persona", Value: "muffin"}})
		if err != nil {
			log.Fatalln(err)
		}

		defer cur.Close(context.TODO())

		cur.All(context.TODO(), &data)
		ch <- 1
	}()

	// 지식 데이터
	go func() {
		cur, err := databases.Database.Learns.Find(context.TODO(), bson.D{{Key: "command", Value: question}})
		if err != nil {
			log.Fatalln(err)
		}

		defer cur.Close(context.TODO())

		cur.All(context.TODO(), &learnData)
		ch <- 1
	}()

	for range 2 {
		<-ch
	}
	close(ch)

	if x > 2 && len(learnData) != 0 {
		data := learnData[rand.Intn(len(learnData))]
		user, _ := s.User(data.UserId)

		result =
			fmt.Sprintf("%s\n%s", data.Result, utils.InlineCode(fmt.Sprintf("%s님이 알려주셨어요.", user.Username)))
	} else {
		result = data[rand.Intn(len(data))].Text
	}
	return result
}

func getAIResponse(question string) string {
	result, err := ChatBot.Gemini.Models.GenerateContent(context.TODO(), configs.Config.Chatbot.Gemini.Model, genai.Text(question), ChatBot.config)
	if err != nil {
		ChatBot.Mode = ChatbotDefault
		return "AI에 문제가 생겼ㅇ어요."
	}
	return result.Text()
}

func (c *Chatbot) GetResponse(question string) string {
	switch c.Mode {
	case ChatbotDefault:
		return getDefaultResponse(c.s, question)
	default:
		return ""
	}
}
