package components

import (
	"context"
	"strings"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var DeregisterComponent *commands.Component = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		customId := ctx.Inter.MessageComponentData().CustomID
		if !strings.HasPrefix(customId, utils.DeregisterAgree) && !strings.HasPrefix(customId, utils.DeregisterDisagree) {
			return false
		}

		if ctx.Inter.User.ID != utils.GetDeregisterUserId(customId) {
			return false
		}
		return true
	},
	Run: func(ctx *commands.ComponentContext) error {
		err := ctx.Inter.DeferUpdate()
		if err != nil {
			return err
		}

		customId := ctx.Inter.MessageComponentData().CustomID
		flags := discordgo.MessageFlagsIsComponentsV2

		switch {
		case strings.HasPrefix(customId, utils.DeregisterAgree):
			filter := databases.User{UserId: ctx.Inter.User.ID}
			_, err := databases.GetDatabase().Users.DeleteOne(context.TODO(), filter)
			if err != nil {
				return err
			}

			_, err = databases.GetDatabase().Learns.DeleteMany(context.TODO(), filter)
			if err != nil {
				return err
			}

			_, err = databases.GetDatabase().Memory.DeleteMany(context.TODO(), filter)
			if err != nil {
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
		case strings.HasPrefix(customId, utils.DeregisterDisagree):
			return ctx.Inter.EditReply(&utils.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					utils.GetCanceledContainer(discordgo.TextDisplay{
						Content: "탈퇴를 거부했어요.",
					}),
				},
			})
		}
		return nil
	},
}
