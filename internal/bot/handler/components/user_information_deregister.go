package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/handler/commands"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Component(utils.UserInformationDeregister+"/{user_id}", func(e *handler.ComponentEvent) error {
			if e.User().ID.String() != e.Vars["user_id"] {
				return e.CreateMessage(
					discord.NewMessageCreateV2(builders.MakeHasNoPermissionContainer()).
						WithEphemeral(true),
				)
			}

			return commands.HandleDeregister(&events.ApplicationCommandInteractionCreate{
				GenericEvent:                  e.GenericEvent,
				ApplicationCommandInteraction: discord.ApplicationCommandInteraction{},
				Respond:                       e.Respond,
			})
		})
	})
}
