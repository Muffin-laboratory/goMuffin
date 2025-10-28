package chat

import (
	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/repository"
)

func Settings(inter *builders.InteractionCreate) error {
	settings, err := repository.NewUserSettings(inter.User)
	if err != nil {
		return err
	}

	return builders.NewMessageSender(inter).
		AddComponents(settings.MakeContainer()).
		SetComponentsV2(true).
		SetEphemeral(true).
		Send()
}
