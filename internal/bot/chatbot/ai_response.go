package chatbot

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/disgoorg/disgo/discord"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/genai"
)

func (c *Chatbot) getAIResponse(ctx context.Context, user discord.User, question string, attachments ...discord.Attachment) (string, error) {
	const twelveHours = 43_200

	dbUser, err := repository.GetDatabase().Users.FindByID(ctx, int64(user.ID))
	if err != nil {
		return "살려주ㅅ세요", err
	}

	if _, err := repository.GetDatabase().Chats.FindOne(ctx, query.ChatQueryBuilder().SetUserID(int64(user.ID))); err != nil {
		if err == mongo.ErrNoDocuments {
			if _, err = repository.GetDatabase().Chats.Create(ctx, int64(user.ID), fmt.Sprintf("새로운 채팅 %06d", rand.Intn(999999))); err != nil {
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
			result, err := repository.GetDatabase().Chats.Create(ctx, int64(user.ID), fmt.Sprintf("새로운 채팅 %06d", rand.Intn(999999)))
			if err != nil {
				return "살려주ㅅ세요", err
			}

			dbUser.ChatID = result.ID
		}
	}

	chat, err := c.GetChat(ctx, &user, dbUser.ChatID)
	if err != nil {
		return "살려주ㅅ세요", err
	}

	var parts []genai.Part
	var files []repository.File

	if len(attachments) != 0 {
		genaiFiles, err := getFiles(c.Gemini, attachments)
		if err != nil {
			return "살려주ㅅ세요", err
		}

		for _, file := range genaiFiles {
			parts = append(parts, *genai.NewPartFromFile(*file))
			files = append(files, repository.File{URI: file.URI, MIMEType: file.MIMEType})
		}
	}

	parts = append(parts, *genai.NewPartFromText(question))

	result, err := chat.SendMessage(ctx, parts...)
	if err != nil {
		return "살려주ㅅ세요", err
	}

	resultText := result.Text()
	if _, err = repository.GetDatabase().Memory.Create(ctx, dbUser.ChatID, int64(user.ID), question, resultText, files); err != nil {
		return "살려주ㅅ세요", err
	}

	log.Printf("%s TOKEN: %d", user.ID, result.UsageMetadata.PromptTokenCount)

	return fmt.Sprintf("%s\n`해당 문장은 AI가 생성한 것이에요.`", resultText), nil
}
