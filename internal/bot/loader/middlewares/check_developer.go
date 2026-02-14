package middlewares

import (
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/disgoorg/disgo/handler"
)

func CheckDeveloperMiddleware() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			developerID := configs.GetConfig().Bot.OwnerID
			if e.User().ID != developerID {
				return nil
			}

			return next(e)
		}
	}
}
