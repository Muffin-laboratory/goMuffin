package chat

import (
	"fmt"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/databases"
)

func chatSendErrorMessage(m any) error {
	return builders.NewMessageSender(m).
		AddComponents(builders.MakeErrorContainer(fmt.Sprintf("채팅모드가 %s여야해요.", databases.ModeString(databases.ChattingAIMode)))).
		SetComponentsV2(true).
		SetReply(true).
		SetEphemeral(true).
		Send()
}
