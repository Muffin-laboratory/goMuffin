package middlewares

import (
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func getUserID(customID string) string {
	userID := getID(customID)
	if strings.Contains(userID, ":") {
		userID = strings.Split(userID, ":")[0]
	}

	return userID
}

func getID(customID string) string {
	parts := strings.Split(customID, "/")
	return parts[len(parts)-1]
}

func CheckIDMiddleware() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			inter, ok := e.Interaction.(discord.ComponentInteraction)
			if !ok {
				return next(e)
			}

			userID := getUserID(inter.Data.CustomID())
			if userID != inter.User().ID.String() {
				_, err := e.CreateFollowupMessage(
					discord.NewMessageCreateV2(builders.MakeHasNoPermissionContainer()).
						WithEphemeral(true),
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
