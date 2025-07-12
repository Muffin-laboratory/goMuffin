package components

import (
	"context"
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
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
					utils.GetDeclineContainer(discordgo.TextDisplay{Content: "당신은 해당 권한이 없ㅇ어요."}),
				},
			})
			return false
		}
		return true
	},
	Run: func(ctx *commands.ComponentContext) error {
		i := ctx.Inter

		err := i.DeferUpdate()
		if err != nil {
			return err
		}

		id, itemId := utils.GetDeleteLearnedDataId(i.MessageComponentData().CustomID)
		_, err = databases.Database.Learns.DeleteOne(context.TODO(), databases.Learn{Id: id})
		if err != nil {
			return err
		}

		flags := discordgo.MessageFlagsIsComponentsV2
		return i.EditReply(&utils.InteractionEdit{
			Flags: &flags,
			Components: &[]discordgo.MessageComponent{
				utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("%d번을 삭ㅈ제했어요.", itemId)}),
			},
		})
	},
}
