package chat

import (
	"context"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
)

func Settings(ctx context.Context, inter *builders.CommandCreate) error {
	settings, err := builders.NewUserSettings(ctx, inter.User())
	if err != nil {
		return err
	}

	return builders.NewMessageSender(inter).
		AddComponents(settings.MakeContainer()).
		SetComponentsV2(true).
		SetEphemeral(true).
		Send()
}
