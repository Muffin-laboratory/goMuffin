package chatbot

import (
	"context"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/cache"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/genai"
)

var chats = cache.New[bson.ObjectID, genai.Chat](time.Hour * 12)

func (c *Chatbot) GetChat(ctx context.Context, user *discordgo.User, chatID bson.ObjectID) (*genai.Chat, error) {
	if cache, ok := chats.Get(chatID); ok {
		return &cache, nil
	}

	contents, err := repository.GetDatabase().Memory.Get(ctx, chatID)
	if err != nil {
		return nil, err
	}

	prompt, err := makePrompt(ctx, c.systemPrompt, user)
	if err != nil {
		return nil, err
	}

	chat, err := c.Gemini.Chats.Create(context.TODO(), configs.GetConfig().Chatbot.Gemini.Model, &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(prompt, genai.RoleUser),
		Tools: []*genai.Tool{
			{
				GoogleSearch: &genai.GoogleSearch{},
			},
		},
	}, contents)
	if err != nil {
		return nil, err
	}

	chats.Set(chatID, *chat)

	return chat, nil
}
