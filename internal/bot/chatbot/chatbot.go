package chatbot

import (
	"context"

	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"google.golang.org/genai"
)

type Chatbot struct {
	Gemini       *genai.Client
	systemPrompt string
	s            *bot.Client
}

var instance *Chatbot

func Make(s *bot.Client) error {
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

func (c *Chatbot) GetResponse(ctx context.Context, user discord.User, question string, attachments ...discord.Attachment) (string, error) {
	mode, err := repository.GetDatabase().Users.GetUserChattingMode(ctx, user.ID.String())
	if err != nil {
		return "살려주ㅅ세요", err
	}

	switch mode {
	case repository.ChattingMuffinMode:
		return c.getMuffinResponse(ctx, question)
	default:
		return c.getAIResponse(ctx, user, question, attachments...)
	}
}
