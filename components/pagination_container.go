package components

import (
	"strings"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var PaginationContainerComponent *commands.Component = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		i := ctx.Inter

		if i.MessageComponentData().ComponentType == discordgo.ButtonComponent {
			customId := i.MessageComponentData().CustomID

			isPrev := strings.HasPrefix(customId, utils.PaginationEmbedPrev)
			isNext := strings.HasPrefix(customId, utils.PaginationEmbedNext)
			isSetPage := strings.HasPrefix(customId, utils.PaginationEmbedPages)
			if !isPrev && !isNext && !isSetPage {
				return false
			}

			id := utils.GetPaginationEmbedID(customId)
			userID := utils.GetPaginationEmbedUserID(id)
			if i.Member.User.ID != userID {
				return false
			}

			if utils.GetPaginationContainer(id) == nil {
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
		p := utils.GetPaginationContainer(id)

		if strings.HasPrefix(customID, utils.PaginationEmbedPrev) {
			p.Prev(ctx.Inter)
			return nil
		} else if strings.HasPrefix(customID, utils.PaginationEmbedNext) {
			p.Next(ctx.Inter)
			return nil
		} else {
			p.ShowModal(ctx.Inter)
			return nil
		}
	},
}

func init() {
	commands.GetDiscommand().LoadComponent(PaginationContainerComponent)
}
