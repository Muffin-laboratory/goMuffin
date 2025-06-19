package components

import (
	"context"
	"fmt"
	"log"
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

		if ctx.Inter.User.ID != utils.GetServicesUserId(customId) {
			return false
		}
		return true
	},
	Run: func(ctx *commands.ComponentContext) {
		ctx.Inter.DeferUpdate()
		customId := ctx.Inter.MessageComponentData().CustomID
		flags := discordgo.MessageFlagsIsComponentsV2

		switch {
		case strings.HasPrefix(customId, utils.ServiceAgree):
			_, err := databases.Database.Users.InsertOne(context.TODO(), databases.InsertUser{
				UserId:    ctx.Inter.User.ID,
				CreatedAt: time.Now(),
			})
			if err != nil {
				log.Println(err)
				ctx.Inter.EditReply(&utils.InteractionEdit{
					Flags: &flags,
					Components: &[]discordgo.MessageComponent{
						utils.GetErrorContainer(discordgo.TextDisplay{
							Content: "가입을 하다가 오류가 생겼어요.",
						}),
					},
				})
				return
			}

			ctx.Inter.EditReply(&utils.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					utils.GetSuccessContainer(discordgo.TextDisplay{
						Content: fmt.Sprintf("가입을 했어요. 이제 %s의 모든 기능을 사용할 수 있어요.", ctx.Inter.Session.State.User.Username),
					}),
				},
			})
			return
		case strings.HasPrefix(customId, utils.ServiceDisagree):
			ctx.Inter.EditReply(&utils.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					utils.GetDeclineContainer(discordgo.TextDisplay{
						Content: "가입을 거부했어요.",
					}),
				},
			})
			return
		}
	},
}
