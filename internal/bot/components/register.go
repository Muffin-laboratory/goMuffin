package components

import (
	"fmt"

	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/commands"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

var RegisterComponent *commands.Component = &commands.Component{
	Parse: func(inter *builders.InteractionCreate) bool {
		customID := inter.MessageComponentData().CustomID
		if !strings.HasPrefix(customID, utils.ServiceAgree) && !strings.HasPrefix(customID, utils.ServiceDisagree) {
			return false
		}

		if inter.User.ID != utils.GetServiceUserID(customID) {
			return false
		}
		return true
	},
	Run: func(inter *builders.InteractionCreate) error {
		if err := inter.DeferUpdate(); err != nil {
			return err
		}

		customID := inter.MessageComponentData().CustomID
		flags := discordgo.MessageFlagsIsComponentsV2

		switch {
		case strings.HasPrefix(customID, utils.ServiceAgree):
			if _, err := repository.GetDatabase().Users.Create(inter.User.ID); err != nil {
				return err
			}

			return inter.EditReply(&builders.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					builders.MakeSuccessContainer(fmt.Sprintf("가입을 했어요. 이제 %s의 모든 기능을 사용할 수 있어요.", inter.Session.State.User.Username)),
				},
			})
		case strings.HasPrefix(customID, utils.ServiceDisagree):
			return inter.EditReply(&builders.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					builders.MakeDeclineContainer("가입을 거부했어요."),
				},
			})
		default:
			return nil
		}
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(RegisterComponent)
}
