package chatbot

import (
	"strings"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func ParseResult(content string, s *discordgo.Session, m *discordgo.MessageCreate) string {
	result := content
	userCreatedAt, _ := discordgo.SnowflakeTimestamp(m.Author.ID)

	result = strings.ReplaceAll(result, "{user.name}", m.Author.Username)
	result = strings.ReplaceAll(result, "{user.mention}", m.Author.Mention())
	result = strings.ReplaceAll(result, "{user.globalName}", m.Author.GlobalName)
	result = strings.ReplaceAll(result, "{user.id}", m.Author.ID)
	result = strings.ReplaceAll(result, "{user.createdAt}", utils.Time(&userCreatedAt, utils.RelativeTime))
	result = strings.ReplaceAll(result, "{user.joinedAt}", utils.Time(&m.Member.JoinedAt, utils.RelativeTime))

	result = strings.ReplaceAll(result, "{muffin.version}", configs.MUFFIN_VERSION)
	result = strings.ReplaceAll(result, "{muffin.updatedAt}", utils.Time(configs.UpdatedAt, utils.RelativeTime))
	result = strings.ReplaceAll(result, "{muffin.startedAt}", utils.Time(configs.StartedAt, utils.RelativeTime))
	result = strings.ReplaceAll(result, "{muffin.name}", s.State.User.Username)
	result = strings.ReplaceAll(result, "{muffin.id}", s.State.User.ID)
	return result
}
