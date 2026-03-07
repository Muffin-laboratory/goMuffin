package middlewares

import (
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func ShowPreviewWarningMessage() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			err := next(e)
			if err != nil {
				return err
			}

			if e.Type() == discord.InteractionTypeApplicationCommand {
				_, err := e.CreateFollowupMessage(
					discord.NewMessageCreateV2(
						discord.NewContainer(
							discord.NewTextDisplay("### ⚠️ 경고"),
							discord.NewTextDisplay("- 해당 버전의 머핀봇은 아직 개발단계로, **매우 불안정한** 상태에요. "+
								"따라서 심각한 오류가 발생할 수 있어요."),
							discord.NewTextDisplayf("-# 현재 버전: %s", configs.MuffinVersion),
						),
					).
						WithEphemeral(true),
				)
				return err
			}

			return nil
		}
	}
}
