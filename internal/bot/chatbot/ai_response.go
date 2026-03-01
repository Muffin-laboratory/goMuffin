package chatbot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/disgoorg/disgo/discord"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/genai"
)

func (c *Chatbot) getAIResponse(ctx context.Context, user discord.User, channelID int64, question string, attachments ...discord.Attachment) (string, error) {
	const twelveHours = 43_200
	var chatID any

	dbUser, err := repository.GetDatabase().Users.FindByID(ctx, int64(user.ID))
	if err != nil {
		return "", err
	}

	if !dbUser.ChatPerChannel {
		channelID = 0
		chatID = dbUser.ChatID

		if _, err := repository.GetDatabase().Chats.FindOne(ctx, query.ChatQueryBuilder().SetUserID(int64(user.ID))); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				if _, err = repository.GetDatabase().Chats.Create(ctx, int64(user.ID), dbUser.Prompt, ""); err != nil {
					return "", err
				}
			} else {
				return "", err
			}
		}

		if dbUser.CreateNewChatAfter12Hours {
			timestamp, err := repository.GetDatabase().Memory.GetLastMemoryTimestamp(ctx, dbUser.ChatID)
			if err != nil {
				return "", err
			}

			if time.Now().Unix()-timestamp > twelveHours {
				result, err := repository.GetDatabase().Chats.Create(ctx, int64(user.ID), dbUser.Prompt, "")
				if err != nil {
					return "", err
				}

				dbUser.ChatID = result.ID
			}
		}
	} else {
		chatID = channelID
	}

	chat, err := c.GetChat(ctx, &user, chatID)
	if err != nil {
		return "", err
	}

	var parts []genai.Part
	var files []repository.File

	if len(attachments) != 0 {
		genaiFiles, err := getFiles(ctx, c.Gemini, attachments)
		if err != nil {
			return "", err
		}

		for _, file := range genaiFiles {
			parts = append(parts, *genai.NewPartFromFile(*file))
			files = append(files, repository.File{URI: file.URI, MIMEType: file.MIMEType})
		}
	}

	parts = append(parts, *genai.NewPartFromText(question))

	result, err := chat.SendMessage(ctx, parts...)
	if err != nil {
		return "", err
	}

	resultText := result.Text()
	if _, err = repository.GetDatabase().Memory.Create(ctx, dbUser.ChatID, int64(user.ID), question, resultText, files, channelID); err != nil {
		return "", err
	}

	log.Printf("%s TOKEN: %d", user.ID, result.UsageMetadata.PromptTokenCount)

	return fmt.Sprintf("%s\n`해당 문장은 AI가 생성한 것이에요.`", resultText), nil
}
