package chatbot

import (
	"strings"
	"time"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func ParseResult(content string, s *discordgo.Session, m any) string {
	result := content

	var user *discordgo.User
	var joinedAt *time.Time

	switch m := m.(type) {
	case *utils.MessageCreate:
		user = m.Author
		joinedAt = &m.Member.JoinedAt
	case *utils.InteractionCreate:
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
