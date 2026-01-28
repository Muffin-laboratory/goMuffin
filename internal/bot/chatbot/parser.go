package chatbot

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
	"google.golang.org/genai"
)

func ParseResult(content string, s *discordgo.Session, m any) string {
	result := content

	var user *discordgo.User
	var joinedAt *time.Time

	switch m := m.(type) {
	case *builders.MessageCreate:
		user = m.Author
		joinedAt = &m.Member.JoinedAt
	case *builders.InteractionCreate:
		user = m.Member.User
		joinedAt = &m.Member.JoinedAt
	}

	userCreatedAt, _ := discordgo.SnowflakeTimestamp(user.ID)

	result = strings.ReplaceAll(result, "{user.name}", user.Username)
	result = strings.ReplaceAll(result, "{user.mention}", user.Mention())
	result = strings.ReplaceAll(result, "{user.globalName}", user.GlobalName)
	result = strings.ReplaceAll(result, "{user.id}", user.ID)
	result = strings.ReplaceAll(result, "{user.createdAt}", utils.Time(&userCreatedAt, utils.RelativeTime))
	result = strings.ReplaceAll(result, "{user.joinedAt}", utils.Time(joinedAt, utils.RelativeTime))

	result = strings.ReplaceAll(result, "{muffin.version}", configs.MuffinVersion)
	result = strings.ReplaceAll(result, "{muffin.updatedAt}", utils.Time(configs.UpdatedAt, utils.RelativeTime))
	result = strings.ReplaceAll(result, "{muffin.startedAt}", utils.Time(configs.StartedAt, utils.RelativeTime))
	result = strings.ReplaceAll(result, "{muffin.name}", s.State.User.Username)
	result = strings.ReplaceAll(result, "{muffin.id}", s.State.User.ID)
	return result
}

func getFiles(client *genai.Client, attachments *[]*discordgo.MessageAttachment) (files []*genai.File, err error) {
	for _, attachment := range *attachments {
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

		file, err = client.Files.Upload(context.TODO(), resp.Body, &genai.UploadFileConfig{
			MIMEType:    attachment.ContentType,
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
