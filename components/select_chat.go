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
		customID := i.MessageComponentData().CustomID

		if !strings.HasPrefix(customID, utils.SelectChat) {
			return false
		}

		userID := utils.GetChatUserID(customID)
		if i.User.ID != userID {
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

		if err := i.DeferUpdate(); err != nil {
			return err
		}

		id, itemID := utils.GetSelectChatID(i.MessageComponentData().CustomID)
		if _, err := databases.GetDatabase().Users.UpdateOne(context.TODO(), databases.User{UserID: i.User.ID}, bson.D{{
			Key:   "$set",
			Value: databases.User{ChatID: id},
		}}); err != nil {
			return err
		}

		flags := discordgo.MessageFlagsIsComponentsV2
		return i.EditReply(&utils.InteractionEdit{
			Flags: &flags,
			Components: &[]discordgo.MessageComponent{
				utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("%d번으로 채팅을 변경했어요.", itemID)}),
			},
		})
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(SelectChatComponent)
}
