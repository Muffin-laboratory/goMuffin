package components

import (
	"fmt"
	"strings"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/commands"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

var DeleteKnowledgeComponent *commands.Component = &commands.Component{
	Parse: func(inter *builders.InteractionCreate) bool {
		customID := inter.MessageComponentData().CustomID

		if !strings.HasPrefix(customID, utils.DeleteKnowledge) {
			return false
		}

		userID := utils.GetDeleteKnowledgeUserID(customID)
		if inter.Member.User.ID != userID {
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

		id, itemID := utils.GetDeleteKnowledgeID(inter.MessageComponentData().CustomID)
		if _, err := repository.GetDatabase().Knowledge.Delete(id); err != nil {
			return err
		}

		flags := discordgo.MessageFlagsIsComponentsV2
		return inter.EditReply(&builders.InteractionEdit{
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
