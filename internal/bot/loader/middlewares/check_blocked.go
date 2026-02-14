package middlewares

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/handler"
)

func CheckBlockedMiddleware() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			blocked, reason := repository.GetDatabase().Users.IsUserBlocked(e.Ctx, int64(e.User().ID))
			if blocked {
				return builders.NewMessageSender(e).
					AddComponents(builders.MakeUserIsBlockedContainer(*e.User().GlobalName, reason)).
					SetComponentsV2(true).
					SetReply(true).
					Send()
			}

			return next(e)
		}
	}
}
