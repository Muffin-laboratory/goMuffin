package middlewares

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/handler"
)

func CheckUserMiddleware() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			if !repository.GetDatabase().Users.IsUser(e.Ctx, int64(e.User().ID)) {
				return builders.NewMessageSender(e).
					AddComponents(builders.MakeUserIsNotRegisteredErrContainer()).
					SetComponentsV2(true).
					SetEphemeral(true).
					SetReply(true).
					Send()
			}

			return next(e)
		}
	}
}

func CheckUserAndBlockedMiddleware() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			next = CheckUserMiddleware()(next)
			next = CheckBlockedMiddleware()(next)
			return next(e)
		}
	}
}
