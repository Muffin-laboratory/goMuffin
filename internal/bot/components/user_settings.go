package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(
			middlewares.CheckUserSettingsMiddleware(),
			middlewares.TimeoutAndDeferMiddleware(loader.Timeout(), discord.InteractionTypeComponent, true, false),
		)

		r.Component(utils.UserSettings+"/{type}/{id}", func(e *handler.ComponentEvent) error {
			settingsType := e.Vars["type"]
			id := e.Vars["id"]
			settings := builders.GetUserSettings(id)

			switch settingsType {
			case "chatting_mode":
				settings.ToggleChattingMode()
			case "reply_user":
				settings.ToggleReplyUser()
			case "12hours":
				settings.Toggle12Hours()
			case "submit":
				err := settings.Submit(e.Ctx)
				if err != nil {
					return err
				}

				_, err = e.UpdateInteractionResponse(
					discord.NewMessageUpdateV2([]discord.LayoutComponent{
						builders.MakeSuccessContainer("- 봇의 대화 설정을 성공적으로 바꾸었어요."),
					}),
				)
				return err
			}

			_, err := e.UpdateInteractionResponse(
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					settings.MakeContainer(),
				}),
			)
			return err
		})
	})
}
