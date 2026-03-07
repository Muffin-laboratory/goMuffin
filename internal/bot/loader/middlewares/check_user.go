package middlewares

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func CheckIsUser() handler.Middleware {
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

func CheckIsUserAndBlocked() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			next = CheckIsUser()(next)
			next = CheckBlocked()(next)
			return next(e)
		}
	}
}
