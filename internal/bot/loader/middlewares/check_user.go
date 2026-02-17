package middlewares

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func CheckUserMiddleware() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			if !repository.GetDatabase().Users.IsUser(e.Ctx, int64(e.User().ID)) {
				return e.CreateMessage(
					discord.NewMessageCreateV2(builders.MakeUserIsNotRegisteredErrContainer()).
						WithEphemeral(true),
				)
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
