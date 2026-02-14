package middlewares

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func CheckPaginationContainerMiddleware() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			next = CheckIDMiddleware()(next)
			if next == nil {
				return nil
			}

			var customID string

			inter, ok := e.Interaction.(discord.ComponentInteraction)
			if !ok {
				inter, ok := e.Interaction.(discord.ModalSubmitInteraction)
				if !ok {
					return nil
				}

				customID = inter.Data.CustomID
			} else {
				customID = inter.Data.CustomID()
			}

			id := getID(customID)
			if builders.GetPaginationContainer(id) == nil {
				return nil
			}

			return next(e)
		}
	}
}
