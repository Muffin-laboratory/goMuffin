package components

import (
	"context"
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var DeleteLearnedDataComponent *commands.Component = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		i := ctx.Inter
		customId := i.MessageComponentData().CustomID

		if !strings.HasPrefix(customId, utils.DeleteLearnedData) {
			return false
		}

		userId := utils.GetDeleteLearnedDataUserId(customId)
		if i.Member.User.ID != userId {
			i.Reply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsIsComponentsV2,
				Components: []discordgo.MessageComponent{
					utils.GetErrorContainer(discordgo.TextDisplay{Content: "당신은 해당 권한이 없ㅇ어요."}),
				},
			},
			)
			return false
		}
		return true
	},
	Run: func(ctx *commands.ComponentContext) {
		i := ctx.Inter

		i.DeferUpdate()

		id, itemId := utils.GetDeleteLearnedDataId(i.MessageComponentData().CustomID)
		fmt.Println(id, itemId)

		databases.Database.Learns.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: id}})

		flags := discordgo.MessageFlagsIsComponentsV2
		i.EditReply(&utils.InteractionEdit{
			Flags: &flags,
			Components: &[]discordgo.MessageComponent{
				utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("%d번을 삭ㅈ제했어요.", itemId)}),
			},
		})
	},
}
