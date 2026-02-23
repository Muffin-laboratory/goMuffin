package middlewares

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func SendErrorMessage() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			err := next(e)
			if err != nil {
				owner, _ := e.Client().Rest.GetUser(configs.Configs().Bot.OwnerID)
				container := builders.MakeErrorContainer("오류가 발생하였어요. 만약 계속 발생한다면, `%s`으로 연락해주세요.\n"+
					"-# 현재 버전: %s", owner.Username, configs.MuffinVersion)
				if err := e.CreateMessage(
					discord.NewMessageCreateV2(container).
						WithEphemeral(true),
				); err != nil {
					e.UpdateInteractionResponse(
						discord.NewMessageUpdateV2([]discord.LayoutComponent{
							container,
						}),
					)
				}
			}

			return err
		}
	}
}
