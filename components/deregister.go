package components

import (
	"strings"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var DeregisterComponent *commands.Component = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		customID := ctx.Inter.MessageComponentData().CustomID
		if !strings.HasPrefix(customID, utils.DeregisterAgree) && !strings.HasPrefix(customID, utils.DeregisterDisagree) {
			return false
		}

		if ctx.Inter.User.ID != utils.GetDeregisterUserID(customID) {
			return false
		}
		return true
	},
	Run: func(ctx *commands.ComponentContext) error {
		err := ctx.Inter.DeferUpdate()
		if err != nil {
			return err
		}

		customID := ctx.Inter.MessageComponentData().CustomID
		flags := discordgo.MessageFlagsIsComponentsV2

		switch {
		case strings.HasPrefix(customID, utils.DeregisterAgree):
			userID := ctx.Inter.User.ID

			if _, err := databases.GetDatabase().Users.Delete(userID); err != nil {
				return err
			}

			if _, err := databases.GetDatabase().Knowledge.DeleteByUserID(userID); err != nil {
				return err
			}

			if _, err := databases.GetDatabase().Memory.DeleteByUserID(userID); err != nil {
				return err
			}

			return ctx.Inter.EditReply(&utils.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					utils.GetSuccessContainer(discordgo.TextDisplay{
						Content: "탈퇴를 했어요.",
					}),
				},
			})
		case strings.HasPrefix(customID, utils.DeregisterDisagree):
			return ctx.Inter.EditReply(&utils.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					utils.GetCanceledContainer(discordgo.TextDisplay{
						Content: "탈퇴를 거부했어요.",
					}),
				},
			})
		default:
			return nil
		}
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(DeregisterComponent)
}
