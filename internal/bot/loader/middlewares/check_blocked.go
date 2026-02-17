package middlewares

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func CheckBlockedMiddleware() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			blocked, reason := repository.GetDatabase().Users.IsUserBlocked(e.Ctx, int64(e.User().ID))
			if blocked {
				return e.CreateMessage(
					discord.NewMessageCreateV2(builders.MakeUserIsBlockedContainer(*e.User().GlobalName, reason)).
						WithEphemeral(true),
				)
			}

			return next(e)
		}
	}
}
