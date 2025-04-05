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
		var userId string
		i := ctx.Inter
		customId := i.MessageComponentData().CustomID

		if i.MessageComponentData().ComponentType == discordgo.ButtonComponent {
			if !strings.HasPrefix(customId, utils.DeleteLearnedDataCancel) {
				return false
			}

			userId = customId[len(utils.DeleteLearnedDataCancel):]
			if i.Member.User.ID == userId {
				i.Update(&discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "❌ 취소",
							Description: "지식 삭제 작업ㅇ을 취소했어요.",
							Color:       utils.EmbedFail,
						},
					},
				})
				return false
			}
		} else {
			if !strings.HasPrefix(customId, utils.DeleteLearnedDataUserId) {
				return false
			}

			userId = customId[len(utils.DeleteLearnedDataUserId):]
		}

		if i.Member.User.ID != userId {
			i.Reply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "❌ 오류",
						Description: "당신은 해당 권한이 없ㅇ어요.",
						Color:       utils.EmbedFail,
					},
				},
				Components: []discordgo.MessageComponent{},
			},
			)
			return false
		}
		return true
	},
	Run: func(ctx *commands.ComponentContext) {
		i := ctx.Inter

		i.DeferUpdate()

		id, _ := bson.ObjectIDFromHex(strings.ReplaceAll(utils.ItemIdRegexp.ReplaceAllString(i.MessageComponentData().Values[0][len(utils.DeleteLearnedData):], ""), "&", ""))
		itemId := strings.ReplaceAll(utils.ItemIdRegexp.FindAllString(i.MessageComponentData().Values[0], 1)[0], "No.", "")

		databases.Learns.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: id}})

		i.EditReply(&discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "✅ 삭제 완료",
					Description: fmt.Sprintf("%s번을 삭ㅈ제했어요.", itemId),
					Color:       utils.EmbedSuccess,
				},
			},
			Components: &[]discordgo.MessageComponent{},
		})
	},
}
