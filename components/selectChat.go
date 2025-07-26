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

var SelectChatComponent *commands.Component = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		i := ctx.Inter
		customId := i.MessageComponentData().CustomID

		if !strings.HasPrefix(customId, utils.SelectChat) {
			return false
		}

		userId := utils.GetChatUserId(customId)
		if i.User.ID != userId {
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

		id, itemId := utils.GetSelectChatId(i.MessageComponentData().CustomID)
		_, err = databases.GetDatabase().Users.UpdateOne(context.TODO(), databases.User{UserId: i.User.ID}, bson.D{{
			Key:   "$set",
			Value: databases.User{ChatId: id},
		}})
		if err != nil {
			return err
		}

		flags := discordgo.MessageFlagsIsComponentsV2
		return i.EditReply(&utils.InteractionEdit{
			Flags: &flags,
			Components: &[]discordgo.MessageComponent{
				utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("%d번으로 채팅을 변경했어요.", itemId)}),
			},
		})
	},
}
