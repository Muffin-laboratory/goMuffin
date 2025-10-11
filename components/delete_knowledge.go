package components

import (
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var DeleteKnowledgeComponent *commands.Component = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		i := ctx.Inter
		customID := i.MessageComponentData().CustomID

		if !strings.HasPrefix(customID, utils.DeleteKnowledge) {
			return false
		}

		userID := utils.GetDeleteKnowledgeUserID(customID)
		if i.Member.User.ID != userID {
			i.Reply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsIsComponentsV2,
				Components: []discordgo.MessageComponent{
					builders.MakeDeclineContainer("당신은 해당 권한이 없ㅇ어요."),
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

		id, itemID := utils.GetDeleteKnowledgeID(i.MessageComponentData().CustomID)
		if _, err := databases.GetDatabase().Knowledge.Delete(id); err != nil {
			return err
		}

		flags := discordgo.MessageFlagsIsComponentsV2
		return i.EditReply(&builders.InteractionEdit{
			Flags: &flags,
			Components: &[]discordgo.MessageComponent{
				builders.MakeSuccessContainer(fmt.Sprintf("%d번을 삭제했어요.", itemID)),
			},
		})
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(DeleteKnowledgeComponent)
}
