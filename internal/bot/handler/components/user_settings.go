package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders/customid"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(
			middlewares.CheckUserSettings(),
		)

		r.Route(customid.UserSettings, func(r handler.Router) {
			r.Component("/chatting_mode/{id}", func(e *handler.ComponentEvent) error {
				settings := builders.GetUserSettings(e.Vars["id"])

				settings.ToggleChattingMode()
				return e.UpdateMessage(
					discord.NewMessageUpdateV2([]discord.LayoutComponent{
						settings.MakeContainer(),
					}),
				)
			})

			r.Component("/reply_user/{id}", func(e *handler.ComponentEvent) error {
				settings := builders.GetUserSettings(e.Vars["id"])

				settings.ToggleReplyUser()
				return e.UpdateMessage(
					discord.NewMessageUpdateV2([]discord.LayoutComponent{
						settings.MakeContainer(),
					}),
				)
			})

			r.Component("/12hours/{id}", func(e *handler.ComponentEvent) error {
				settings := builders.GetUserSettings(e.Vars["id"])

				settings.Toggle12Hours()
				return e.UpdateMessage(
					discord.NewMessageUpdateV2([]discord.LayoutComponent{
						settings.MakeContainer(),
					}),
				)
			})

			r.Component("/prompt/{id}", func(e *handler.ComponentEvent) error {
				return builders.GetUserSettings(e.Vars["id"]).PromptModal(e)
			})

			r.Component("/chat_per_channel/{id}", func(e *handler.ComponentEvent) error {
				settings := builders.GetUserSettings(e.Vars["id"])

				settings.ToggleChatPerChannel()
				return e.UpdateMessage(
					discord.NewMessageUpdateV2([]discord.LayoutComponent{
						settings.MakeContainer(),
					}),
				)
			})

			r.Component("/cancel/{id}", func(e *handler.ComponentEvent) error {
				builders.GetUserSettings(e.Vars["id"]).Cancel()

				return e.UpdateMessage(
					discord.NewMessageUpdateV2([]discord.LayoutComponent{
						builders.MakeCanceledContainer("- 해당 변경사항을 취소했어요."),
					}),
				)
			})

			r.Group(func(r handler.Router) {
				r.Use(middlewares.TimeoutAndDefer(loader.Timeout(), discord.InteractionTypeComponent, true, false))

				r.Component("/submit/{id}", func(e *handler.ComponentEvent) error {
					err := builders.GetUserSettings(e.Vars["id"]).Submit(e.Ctx)
					if err != nil {
						return err
					}

					_, err = e.UpdateInteractionResponse(
						discord.NewMessageUpdateV2([]discord.LayoutComponent{
							builders.MakeSuccessContainer("- 봇의 대화 설정을 성공적으로 바꾸었어요."),
						}),
					)
					return err
				})
			})
		})
	})
}
