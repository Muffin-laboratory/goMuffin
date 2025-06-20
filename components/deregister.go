package components

import (
	"context"
	"log"
	"strings"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	Run: func(ctx *commands.ComponentContext) {
		ctx.Inter.DeferUpdate()
		customId := ctx.Inter.MessageComponentData().CustomID
		flags := discordgo.MessageFlagsIsComponentsV2

		switch {
		case strings.HasPrefix(customId, utils.DeregisterAgree):
			filter := bson.D{{Key: "user_id", Value: ctx.Inter.User.ID}}
			_, err := databases.Database.Users.DeleteOne(context.TODO(), filter)
			if err != nil {
				log.Println(err)
				// 나중에 에러처리 바꿔야지 :(
				return
			}

			_, err = databases.Database.Learns.DeleteMany(context.TODO(), filter)
			if err != nil {
				log.Println(err)
				// 나중에 에러처리 바꿔야지 :(
				return
			}

			_, err = databases.Database.Memory.DeleteMany(context.TODO(), filter)
			if err != nil {
				log.Println(err)
				// 나중에 에러처리 바꿔야지 :(
				return
			}

			ctx.Inter.EditReply(&utils.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					utils.GetSuccessContainer(discordgo.TextDisplay{
						Content: "탈퇴를 했어요.",
					}),
				},
			})
			return
		case strings.HasPrefix(customId, utils.DeregisterDisagree):
			ctx.Inter.EditReply(&utils.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					utils.GetDeclineContainer(discordgo.TextDisplay{
						Content: "탈퇴를 거부했어요.",
					}),
				},
			})
			return
		}
	},
}
