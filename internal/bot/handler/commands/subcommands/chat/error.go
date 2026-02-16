package chat

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func chatSendErrorMessage(e *handler.CommandEvent) error {
	_, err := e.UpdateInteractionResponse(
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			builders.MakeErrorContainer("채팅모드가 %s여야해요.", repository.ModeString(repository.ChattingAIMode)),
		}),
	)
	return err
}
