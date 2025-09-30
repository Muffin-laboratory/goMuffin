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
		customID := i.MessageComponentData().CustomID

		if !strings.HasPrefix(customID, utils.DeleteChat) && !strings.HasPrefix(customID, utils.DeleteChatCancel) {
			return false
		}

		userID := utils.GetChatUserID(customID)
		if i.Member.User.ID != userID {
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

		if err := i.DeferUpdate(); err != nil {
			return err
		}

		id, name := utils.GetChatID(i.MessageComponentData().CustomID)

		if _, err := databases.GetDatabase().Chats.DeleteOne(context.TODO(), databases.Chat{ID: id}); err != nil {
			return err
		}

		if _, err := databases.GetDatabase().Memory.DeleteByChatID(id); err != nil {
			return err
		}

		flags := discordgo.MessageFlagsIsComponentsV2
		return i.EditReply(&utils.InteractionEdit{
			Flags: &flags,
			Components: &[]discordgo.MessageComponent{
				utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("`%s`번을 삭제했어요.", name)}),
			},
		})
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(DeleteChatComponent)
}
