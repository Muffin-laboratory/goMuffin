package chat

import (
	"fmt"

	"github.com/Muffin-laboratory/goMuffin/builders"
	"github.com/Muffin-laboratory/goMuffin/repository"
)

func chatSendErrorMessage(m any) error {
	return builders.NewMessageSender(m).
		AddComponents(builders.MakeErrorContainer(fmt.Sprintf("채팅모드가 %s여야해요.", repository.ModeString(repository.ChattingAIMode)))).
		SetComponentsV2(true).
		SetReply(true).
		SetEphemeral(true).
		Send()
}
