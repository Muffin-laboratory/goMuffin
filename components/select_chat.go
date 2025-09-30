package components

import (
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
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

		id, name := utils.GetChatID(i.MessageComponentData().CustomID)

		if _, err := databases.GetDatabase().Users.Update(i.User.ID, databases.User{
			ChatID: id,
		}); err != nil {
			return err
		}

		flags := discordgo.MessageFlagsIsComponentsV2
		return i.EditReply(&utils.InteractionEdit{
			Flags: &flags,
			Components: &[]discordgo.MessageComponent{
				utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("`%s`으로 채팅을 변경했어요.", name)}),
			},
		})
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(SelectChatComponent)
}
