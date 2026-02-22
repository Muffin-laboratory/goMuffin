package chatbot

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"google.golang.org/genai"
)

func ParseResult(content string, m any) string {
	result := content

	var user discord.User
	var bot discord.OAuth2User
	var joinedAt *time.Time

	switch m := m.(type) {
	case *events.MessageCreate:
		bot, _ = m.Client().Caches.SelfUser()
		user = m.Message.Author
		joinedAt = m.Message.Member.JoinedAt
	case *handler.CommandEvent:
		user = m.User()

		joinedAt = m.Member().JoinedAt
	}

	userCreatedAt := user.ID.Time()

	result = strings.ReplaceAll(result, "{user.name}", user.Username)
	result = strings.ReplaceAll(result, "{user.mention}", user.Mention())
	result = strings.ReplaceAll(result, "{user.globalName}", *user.GlobalName)
	result = strings.ReplaceAll(result, "{user.id}", user.ID.String())
	result = strings.ReplaceAll(result, "{user.createdAt}", builders.Time(&userCreatedAt, builders.RelativeTime))
	result = strings.ReplaceAll(result, "{user.joinedAt}", builders.Time(joinedAt, builders.RelativeTime))

	result = strings.ReplaceAll(result, "{muffin.version}", configs.MuffinVersion)
	result = strings.ReplaceAll(result, "{muffin.updatedAt}", builders.Time(configs.UpdatedAt(), builders.RelativeTime))
	result = strings.ReplaceAll(result, "{muffin.startedAt}", builders.Time(configs.StartedAt, builders.RelativeTime))
	result = strings.ReplaceAll(result, "{muffin.name}", bot.Username)
	result = strings.ReplaceAll(result, "{muffin.id}", bot.ID.String())
	return result
}

func getFiles(ctx context.Context, client *genai.Client, attachments []discord.Attachment) (files []*genai.File, err error) {
	for _, attachment := range attachments {
		var file *genai.File
		var resp *http.Response

		resp, err = http.Get(attachment.URL)
		if err != nil {
			break
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			break
		}

		file, err = client.Files.Upload(ctx, resp.Body, &genai.UploadFileConfig{
			MIMEType:    *attachment.ContentType,
			DisplayName: attachment.Filename,
		})

		resp.Body.Close()

		if err != nil {
			break
		}

		files = append(files, file)
	}

	return
}
