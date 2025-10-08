package components

import (
	"strings"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var PaginationContainerComponent *commands.Component = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		i := ctx.Inter

		if i.MessageComponentData().ComponentType == discordgo.ButtonComponent {
			customID := i.MessageComponentData().CustomID

			isPrev := strings.HasPrefix(customID, utils.PaginationEmbedPrev)
			isNext := strings.HasPrefix(customID, utils.PaginationEmbedNext)
			isSetPage := strings.HasPrefix(customID, utils.PaginationEmbedPages)
			if !isPrev && !isNext && !isSetPage {
				return false
			}

			id := utils.GetPaginationEmbedID(customID)
			userID := utils.GetPaginationEmbedUserID(id)
			if i.Member.User.ID != userID {
				return false
			}

			if builders.GetPaginationContainer(id) == nil {
				return false
			}
		} else {
			return false
		}
		return true
	},
	Run: func(ctx *commands.ComponentContext) error {
		customID := ctx.Inter.MessageComponentData().CustomID
		id := utils.GetPaginationEmbedID(customID)
		p := builders.GetPaginationContainer(id)

		if strings.HasPrefix(customID, utils.PaginationEmbedPrev) {
			return p.Prev(ctx.Inter)
		} else if strings.HasPrefix(customID, utils.PaginationEmbedNext) {
			return p.Next(ctx.Inter)
		} else {
			return p.ShowModal(ctx.Inter)
		}
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(PaginationContainerComponent)
}
