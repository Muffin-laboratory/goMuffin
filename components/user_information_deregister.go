package components

import (
	"strings"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var UserInformationDeregisterComponent = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		i := ctx.Inter
		customID := i.MessageComponentData().CustomID

		if !strings.HasPrefix(customID, utils.UserInformationDeregister) {
			return false
		}

		userID := utils.GetUserInformationDeregisterUserID(customID)
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
		return commands.DeregisterCommand.Run(&commands.ChatInputContext{
			Inter:   ctx.Inter,
			Command: nil,
		})
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(UserInformationDeregisterComponent)
}
