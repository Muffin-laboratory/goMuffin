package chatbot

import (
	"context"
	"time"

	"github.com/Muffin-laboratory/goMuffin/cache"
	"github.com/Muffin-laboratory/goMuffin/configs"
	"github.com/Muffin-laboratory/goMuffin/repository"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/genai"
)

var chats = cache.New[*genai.Chat](time.Hour * 12)

func (c *Chatbot) GetChat(user *discordgo.User, chatID bson.ObjectID) (*genai.Chat, error) {
	if cache, ok := chats.Get(chatID.Hex()); ok {
		return cache, nil
	}

	contents, err := repository.GetDatabase().Memory.Get(chatID)
	if err != nil {
		return nil, err
	}

	prompt, err := makePrompt(c.systemPrompt, user)
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

	chats.Set(chatID.Hex(), chat)

	return chat, nil
}
