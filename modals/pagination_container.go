package modals

import (
	"strconv"
	"strings"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var PaginationContainerModal *commands.Modal = &commands.Modal{
	Parse: func(inter *builders.InteractionCreate) bool {
		data := inter.ModalSubmitData()
		customID := data.CustomID

		if data.Components[0].Type() != discordgo.ActionsRowComponent {
			return false
		}

		if !strings.HasPrefix(customID, utils.PaginationContainerModal) {
			return false
		}

		id := utils.GetPaginationContainerID(customID)
		userID := utils.GetUserID(id)

		if inter.Member.User.ID != userID {
			return false
		}

		if builders.GetPaginationContainer(id) == nil {
			return false
		}

		cmp := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)

		if _, err := strconv.Atoi(cmp.Value); err != nil {
			inter.Reply(&discordgo.InteractionResponseData{
				Components: []discordgo.MessageComponent{
					builders.MakeErrorContainer("해당 값은 숫자여야해요."),
				},
				Flags: discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsIsComponentsV2,
			})
			return false
		}

		return true
	},
	Run: func(inter *builders.InteractionCreate) error {
		data := inter.ModalSubmitData()
		customID := data.CustomID
		id := utils.GetPaginationContainerID(customID)
		p := builders.GetPaginationContainer(id)
		cmp := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)

		page, _ := strconv.Atoi(cmp.Value)

		return p.Set(inter, page)
	},
}

func init() {
	commands.GetDiscommand().LoadModal(PaginationContainerModal)
}
