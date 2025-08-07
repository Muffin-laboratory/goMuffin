package components

import (
	"context"
	"fmt"

	"strings"
	"time"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var RegisterComponent *commands.Component = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		customId := ctx.Inter.MessageComponentData().CustomID
		if !strings.HasPrefix(customId, utils.ServiceAgree) && !strings.HasPrefix(customId, utils.ServiceDisagree) {
			return false
		}

		if ctx.Inter.User.ID != utils.GetServiceUserId(customId) {
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
		case strings.HasPrefix(customId, utils.ServiceAgree):
			_, err := databases.Database.Users.InsertOne(context.TODO(), databases.User{
				UserId:    ctx.Inter.User.ID,
				CreatedAt: time.Now(),
			})
			if err != nil {
				return err
			}

			return ctx.Inter.EditReply(&utils.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					utils.GetSuccessContainer(discordgo.TextDisplay{
						Content: fmt.Sprintf("가입을 했어요. 이제 %s의 모든 기능을 사용할 수 있어요.", ctx.Inter.Session.State.User.Username),
					}),
				},
			})
		case strings.HasPrefix(customId, utils.ServiceDisagree):
			return ctx.Inter.EditReply(&utils.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					utils.GetDeclineContainer(discordgo.TextDisplay{
						Content: "가입을 거부했어요.",
					}),
				},
			})
		}
		return nil
	},
}
