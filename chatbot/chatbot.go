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
		Mode:   ChatbotAI,
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
	bin, err := os.ReadFile(configs.Config.Chatbot.Gemini.PromptPath)
	if err != nil {
		return err
	}

	ChatBot.config = &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(string(bin), genai.RoleUser),
	}
	return nil
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

func getAIResponse(userId, question string) string {
	contents, err := GetMemory(userId)
	if err != nil {
		ChatBot.Mode = ChatbotMuffin
		log.Fatalln(err)
		return "AI에 문제가 생겼ㅇ어요."
	}

	contents = append(contents, genai.NewContentFromText(question, genai.RoleUser))
	result, err := ChatBot.Gemini.Models.GenerateContent(context.TODO(), configs.Config.Chatbot.Gemini.Model, contents, ChatBot.config)
	if err != nil {
		ChatBot.Mode = ChatbotMuffin
		log.Fatalln(err)
		return "AI에 문제가 생겼ㅇ어요."
	}

	resultText := result.Text()
	err = SaveMemory(&databases.InsertMemory{
		UserId:  userId,
		Content: question,
		Answer:  resultText,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return resultText
}

func (c *Chatbot) GetResponse(userId, question string) string {
	switch c.Mode {
	case ChatbotMuffin:
		return getDefaultResponse(c.s, question)
	default:
		return getAIResponse(userId, question)
	}
}
