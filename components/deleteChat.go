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

var DeleteChatComponent = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		i := ctx.Inter
		customId := i.MessageComponentData().CustomID

		if !strings.HasPrefix(customId, utils.DeleteChat) && !strings.HasPrefix(customId, utils.DeleteChatCancel) {
			return false
		}

		userId := utils.GetChatUserId(customId)
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
		customId := i.MessageComponentData().CustomID

		if strings.HasPrefix(customId, utils.DeleteChatCancel) {
			return i.Update(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsIsComponentsV2,
				Components: []discordgo.MessageComponent{
					utils.GetCanceledContainer(discordgo.TextDisplay{Content: "아무 채팅방을 삭제하지 않았어요."}),
				},
			})
		}

		err := i.DeferUpdate()
		if err != nil {
			return err
		}

		id, itemId := utils.GetDeleteLearnedDataId(i.MessageComponentData().CustomID)
		_, err = databases.GetDatabase().Chats.DeleteOne(context.TODO(), databases.Chat{Id: id})
		if err != nil {
			return err
		}

		_, err = databases.GetDatabase().Memory.DeleteMany(context.TODO(), databases.Memory{ChatId: id})
		if err != nil {
			return err
		}

		flags := discordgo.MessageFlagsIsComponentsV2
		if itemId == 0 {
			return i.EditReply(&utils.InteractionEdit{
				Flags: &flags,
				Components: &[]discordgo.MessageComponent{
					utils.GetSuccessContainer(discordgo.TextDisplay{Content: "해당 채팅을 삭제했어요."}),
				},
			})
		}

		return i.EditReply(&utils.InteractionEdit{
			Flags: &flags,
			Components: &[]discordgo.MessageComponent{
				utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("%d번을 삭제했어요.", itemId)}),
			},
		})
	},
}
