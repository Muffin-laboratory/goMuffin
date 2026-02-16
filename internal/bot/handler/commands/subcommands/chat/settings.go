package chat

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func Settings(e *handler.CommandEvent) error {
	settings, err := builders.NewUserSettings(e.Ctx, e.User())
	if err != nil {
		return err
	}

	_, err = e.UpdateInteractionResponse(
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			settings.MakeContainer(),
		}),
	)
	return err
}
