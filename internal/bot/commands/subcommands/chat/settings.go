package chat

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
)

func Settings(inter *builders.InteractionCreate) error {
	settings, err := repository.NewUserSettings(inter.Ctx, inter.User)
	if err != nil {
		return err
	}

	return builders.NewMessageSender(inter).
		AddComponents(settings.MakeContainer()).
		SetComponentsV2(true).
		SetEphemeral(true).
		Send()
}
