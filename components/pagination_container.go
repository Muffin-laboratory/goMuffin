package components

import (
	"strings"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/utils"
)

var PaginationContainerComponent *commands.Component = &commands.Component{
	Parse: func(inter *builders.InteractionCreate) bool {
		customID := inter.MessageComponentData().CustomID

		isPrev := strings.HasPrefix(customID, utils.PaginationContainerPrev)
		isNext := strings.HasPrefix(customID, utils.PaginationContainerNext)
		isSetPage := strings.HasPrefix(customID, utils.PaginationContainerPages)
		if !isPrev && !isNext && !isSetPage {
			return false
		}

		id := utils.GetPaginationContainerID(customID)
		if inter.Member.User.ID != utils.GetUserID(id) {
			return false
		}

		return builders.GetPaginationContainer(id) != nil
	},
	Run: func(inter *builders.InteractionCreate) error {
		customID := inter.MessageComponentData().CustomID
		id := utils.GetPaginationContainerID(customID)
		p := builders.GetPaginationContainer(id)

		switch {
		case strings.HasPrefix(customID, utils.PaginationContainerPrev):
			return p.Prev(inter)
		case strings.HasPrefix(customID, utils.PaginationContainerNext):
			return p.Next(inter)
		case strings.HasPrefix(customID, utils.PaginationContainerModal):
			return p.ShowModal(inter)
		default:
			return nil
		}

	},
}

func init() {
	commands.GetDiscommand().LoadComponent(PaginationContainerComponent)
}
