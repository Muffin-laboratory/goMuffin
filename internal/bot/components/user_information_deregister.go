package components

import (
	"context"
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/commands"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var UserInformationDeregisterComponent = &loader.Component{
	Parse: func(ctx context.Context, inter *events.ComponentInteractionCreate) bool {
		customID := inter.Data.CustomID()

		if !strings.HasPrefix(customID, utils.UserInformationDeregister) {
			return false
		}

		userID := utils.GetUserInformationDeregisterUserID(customID)
		if inter.User().ID.String() != userID {
			inter.CreateMessage(
				discord.NewMessageCreateBuilder().
					SetComponents(builders.MakeDeclineContainer("당신은 해당 권한이 없ㅇ어요.")).
					SetIsComponentsV2(true).
					SetEphemeral(true).
					Build(),
			)
			return false
		}
		return true
	},
	Run: func(ctx context.Context, inter *events.ComponentInteractionCreate) error {
		return commands.DeregisterCommand.Run(ctx, &builders.CommandCreate{
			ApplicationCommandInteractionCreate: &events.ApplicationCommandInteractionCreate{
				GenericEvent:                  inter.GenericEvent,
				ApplicationCommandInteraction: discord.ApplicationCommandInteraction{},
				Respond:                       inter.Respond,
			}})
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(UserInformationDeregisterComponent)
}
