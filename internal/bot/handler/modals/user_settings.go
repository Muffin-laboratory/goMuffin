package modals

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
		r.Use(middlewares.CheckUserSettingsMiddleware())

		r.Modal(customid.UserSettingsPrompt+"/{id}", func(e *handler.ModalEvent) error {
			var prompt string

			settings := builders.GetUserSettings(e.Vars["id"])

			if opt, ok := e.Data.OptText(customid.UserSettingsPromptSet); ok {
				prompt = opt
			}

			settings.SetPrompt(prompt)

			return e.UpdateMessage(
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					settings.MakeContainer(),
				}),
			)
		})
	})
}
