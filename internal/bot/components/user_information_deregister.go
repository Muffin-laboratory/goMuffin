package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/commands"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

var UserInformationDeregisterComponent = &loader.Component{
	Handle: func(r handler.Router) {
		r.Component(utils.UserInformationDeregister+"/{user_id}", func(inter *handler.ComponentEvent) error {
			if inter.User().ID.String() != inter.Vars["user_id"] {
				return inter.CreateMessage(
					discord.NewMessageCreateBuilder().
						SetComponents(builders.MakeHasNoPermissionContainer()).
						SetIsComponentsV2(true).
						SetEphemeral(true).
						Build(),
				)
			}

			return commands.DeregisterCommand.Run(inter.Ctx, &builders.CommandCreate{
				ApplicationCommandInteractionCreate: &events.ApplicationCommandInteractionCreate{
					GenericEvent:                  inter.GenericEvent,
					ApplicationCommandInteraction: discord.ApplicationCommandInteraction{},
					Respond:                       inter.Respond,
				}})
		})
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(UserInformationDeregisterComponent)
}
