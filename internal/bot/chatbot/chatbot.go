package chatbot

import (
	"context"
	"log/slog"

	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"google.golang.org/genai"
)

type Chatbot struct {
	Gemini       *genai.Client
	systemPrompt string
	corePrompt   string
	s            *bot.Client
}

var instance *Chatbot

func Make(s *bot.Client) error {
	gemini, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  configs.Configs().Chatbot.Gemini.Token,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return err
	}

	instance = &Chatbot{
		Gemini: gemini,
		s:      s,
	}

	system, core, err := loadPrompt()
	if err != nil {
		return err
	}

	instance.systemPrompt = system
	instance.corePrompt = core
	slog.Info("chatbot is created.")
	return nil
}

func GetChatBot() *Chatbot {
	return instance
}

func (c *Chatbot) GetPrompt() string {
	return c.systemPrompt
}

func (c *Chatbot) GetResponse(ctx context.Context, user discord.User, channelID snowflake.ID, question string, attachments ...discord.Attachment) (string, error) {
	mode, err := repository.GetDatabase().Users.GetUserChattingMode(ctx, int64(user.ID))
	if err != nil {
		return "살려주ㅅ세요", err
	}

	switch mode {
	case repository.ChattingMuffinMode:
		return c.getMuffinResponse(ctx, question)
	default:
		return c.getAIResponse(ctx, user, int64(channelID), question, attachments...)
	}
}
