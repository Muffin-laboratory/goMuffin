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

var chats = cache.New[any, *genai.Chat](time.Hour * 12)

func (c *Chatbot) GetChat(ctx context.Context, user *discord.User, chatID any) (*genai.Chat, error) {
	if cache, ok := chats.Get(chatID); ok {
		return cache, nil
	}

	var isChat bool

	filter := query.MemoryQueryBuilder()
	switch chatID := chatID.(type) {
	case int64:
		filter.SetChannelID(chatID)
	case bson.ObjectID:
		filter.SetChatID(chatID)
		isChat = true
	}

	memory, err := repository.GetDatabase().Memory.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	systemPrompt := c.systemPrompt

	if isChat {
		dbChat, err := repository.GetDatabase().Chats.FindByID(ctx, chatID.(bson.ObjectID))
		if err != nil {
			return nil, err
		}

		if dbChat.Prompt != "" {
			systemPrompt = dbChat.Prompt
		}
	}

	prompt, err := makePrompt(ctx, systemPrompt, c.corePrompt, user, isChat)
	if err != nil {
		return nil, err
	}

	var content []*genai.Content
	for _, memory := range memory {
		content = append(content, memory.ToContents()...)
	}

	chat, err := c.Gemini.Chats.Create(ctx, configs.Configs().Chatbot.Gemini.Model, &genai.GenerateContentConfig{
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
