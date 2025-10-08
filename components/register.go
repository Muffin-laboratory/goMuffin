package components

import (
	"fmt"

	"strings"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var RegisterComponent *commands.Component = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		customID := ctx.Inter.MessageComponentData().CustomID
		if !strings.HasPrefix(customID, utils.ServiceAgree) && !strings.HasPrefix(customID, utils.ServiceDisagree) {
			return false
		}

		if ctx.Inter.User.ID != utils.GetServiceUserID(customID) {
			return false
		}
		return true
	},
	Run: func(ctx *commands.ComponentContext) error {
		if err := ctx.Inter.DeferUpdate(); err != nil {
			return err
		}

		customID := ctx.Inter.MessageComponentData().CustomID
		flags := discordgo.MessageFlagsIsComponentsV2

		switch {
		case strings.HasPrefix(customID, utils.ServiceAgree):
			if _, err := databases.GetDatabase().Users.Create(ctx.Inter.User.ID); err != nil {
				return err
			}

			return ctx.Inter.EditReply(&builders.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					builders.MakeSuccessContainer(fmt.Sprintf("가입을 했어요. 이제 %s의 모든 기능을 사용할 수 있어요.", ctx.Inter.Session.State.User.Username)),
				},
			})
		case strings.HasPrefix(customID, utils.ServiceDisagree):
			return ctx.Inter.EditReply(&builders.InteractionEdit{
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
