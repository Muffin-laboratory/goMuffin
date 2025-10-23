package components

import (
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/repository"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var SelectChatComponent *commands.Component = &commands.Component{
	Parse: func(inter *builders.InteractionCreate) bool {
		customID := inter.MessageComponentData().CustomID

		if !strings.HasPrefix(customID, utils.SelectChat) {
			return false
		}

		userID := utils.GetChatUserID(customID)
		if inter.User.ID != userID {
			inter.Reply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsIsComponentsV2,
				Components: []discordgo.MessageComponent{
					builders.MakeDeclineContainer("당신은 해당 권한이 없ㅇ어요."),
				},
			})
			return false
		}
		return true
	},
	Run: func(inter *builders.InteractionCreate) error {
		if err := inter.DeferUpdate(); err != nil {
			return err
		}

		id, name := utils.GetChatID(inter.MessageComponentData().CustomID)

		if _, err := repository.GetDatabase().Users.Update(inter.User.ID, &repository.UserUpdate{
			ChatID: &id,
		}); err != nil {
			return err
		}

		flags := discordgo.MessageFlagsIsComponentsV2
		return inter.EditReply(&builders.InteractionEdit{
			Flags: &flags,
			Components: &[]discordgo.MessageComponent{
				builders.MakeSuccessContainer(fmt.Sprintf("`%s`으로 채팅을 변경했어요.", name)),
			},
		})
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(SelectChatComponent)
}
