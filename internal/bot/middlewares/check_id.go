package middlewares

import (
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func CheckIDMiddleware() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			inter, ok := e.Interaction.(discord.ComponentInteraction)
			if !ok {
				return next(e)
			}

			parts := strings.Split(inter.Data.CustomID(), "/")
			userID := parts[len(parts)-1]
			if userID != inter.User().ID.String() {
				_, err := e.CreateFollowupMessage(
					discord.NewMessageCreateBuilder().
						SetComponents(builders.MakeHasNoPermissionContainer()).
						SetIsComponentsV2(true).
						SetEphemeral(true).
						Build(),
				)
				if err != nil {
					return err
				}

				return nil
			}

			return next(e)
		}
	}
}
