package modals

import (
	"strconv"
	"strings"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var PaginationContainerModal *commands.Modal = &commands.Modal{
	Parse: func(ctx *commands.ModalContext) bool {
		i := ctx.Inter
		data := i.ModalSubmitData()
		customID := data.CustomID

		if data.Components[0].Type() != discordgo.ActionsRowComponent {
			return false
		}

		if !strings.HasPrefix(customID, utils.PaginationEmbedModal) {
			return false
		}

		id := utils.GetPaginationEmbedID(customID)
		userID := utils.GetPaginationEmbedUserID(id)

		if i.Member.User.ID != userID {
			return false
		}

		if utils.GetPaginationContainer(id) == nil {
			return false
		}

		cmp := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)

		if _, err := strconv.Atoi(cmp.Value); err != nil {
			i.Reply(&discordgo.InteractionResponseData{
				Components: []discordgo.MessageComponent{
					utils.GetErrorContainer(discordgo.TextDisplay{Content: "해당 값은 숫자여야해요."}),
				},
				Flags: discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsIsComponentsV2,
			})
			return false
		}

		return true
	},
	Run: func(ctx *commands.ModalContext) error {
		data := ctx.Inter.ModalSubmitData()
		customID := data.CustomID
		id := utils.GetPaginationEmbedID(customID)
		p := utils.GetPaginationContainer(id)
		cmp := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)

		page, _ := strconv.Atoi(cmp.Value)

		return p.Set(ctx.Inter, page)
	},
}
