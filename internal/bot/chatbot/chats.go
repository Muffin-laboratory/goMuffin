package chatbot

import (
	"context"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/cache"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/disgoorg/disgo/discord"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/genai"
)

var chats = cache.New[bson.ObjectID, *genai.Chat](time.Hour * 12)

func (c *Chatbot) GetChat(ctx context.Context, user *discord.User, chatID bson.ObjectID) (*genai.Chat, error) {
	if cache, ok := chats.Get(chatID); ok {
		return cache, nil
	}

	memory, err := repository.GetDatabase().Memory.Find(ctx, query.MemoryQueryBuilder().SetChatID(chatID))
	if err != nil {
		return nil, err
	}

	systemPrompt := c.systemPrompt
	dbChat, err := repository.GetDatabase().Chats.FindByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if dbChat.Prompt != "" {
		systemPrompt = dbChat.Prompt
	}

	prompt, err := makePrompt(ctx, systemPrompt, c.corePrompt, user)
	if err != nil {
		return nil, err
	}

	var content []*genai.Content
	for _, memory := range memory {
		content = append(content, memory.ToContents()...)
	}

	chat, err := c.Gemini.Chats.Create(ctx, configs.GetConfig().Chatbot.Gemini.Model, &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(prompt, genai.RoleUser),
		Tools: []*genai.Tool{
			{
				GoogleSearch: &genai.GoogleSearch{},
			},
		},
	}, content)
	if err != nil {
		return nil, err
	}

	chats.Set(chatID, chat)

	return chat, nil
}
