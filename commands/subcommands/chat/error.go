package chat

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

func chatSendErrorMessage(m any) error {
	return utils.NewMessageSender(m).
		AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: fmt.Sprintf("채팅모드가 %s여야해요.", databases.ModeString(databases.ChattingAIMode))})).
		SetComponentsV2(true).
		SetReply(true).
		SetEphemeral(true).
		Send()
}
