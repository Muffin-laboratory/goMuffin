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

		isFirst := strings.HasPrefix(customID, utils.PaginationContainerFirst)
		isPrev := strings.HasPrefix(customID, utils.PaginationContainerPrev)
		isNext := strings.HasPrefix(customID, utils.PaginationContainerNext)
		isLast := strings.HasPrefix(customID, utils.PaginationContainerLast)
		isSetPage := strings.HasPrefix(customID, utils.PaginationContainerPages)
		if !isFirst && !isPrev && !isNext && !isLast && !isSetPage {
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
		case strings.HasPrefix(customID, utils.PaginationContainerFirst):
			return p.First(inter)
		case strings.HasPrefix(customID, utils.PaginationContainerPrev):
			return p.Prev(inter)
		case strings.HasPrefix(customID, utils.PaginationContainerNext):
			return p.Next(inter)
		case strings.HasPrefix(customID, utils.PaginationContainerLast):
			return p.Last(inter)
		case strings.HasPrefix(customID, utils.PaginationContainerPages):
			return p.ShowModal(inter)
		default:
			return nil
		}

	},
}

func init() {
	commands.GetDiscommand().LoadComponent(PaginationContainerComponent)
}
