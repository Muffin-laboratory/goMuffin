package middlewares

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func CheckIsDeveloper() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			developerID := configs.GetConfig().Bot.OwnerID
			if e.User().ID != developerID {
				return e.CreateMessage(
					discord.NewMessageCreateV2(builders.MakeDeclineContainer("이 명령어는 개발자 전용 명령어에요.")).
						WithEphemeral(true),
				)
			}

			return next(e)
		}
	}
}
