package components

import (
	"strings"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var PaginationEmbedComponent *commands.Component = &commands.Component{
	Parse: func(ctx *commands.ComponentContext) bool {
		i := ctx.Inter

		if i.MessageComponentData().ComponentType == discordgo.ButtonComponent {
			customId := i.MessageComponentData().CustomID

			if !strings.HasPrefix(customId, utils.PaginationEmbedPrev) && !strings.HasPrefix(customId, utils.PaginationEmbedNext) {
				return false
			}

			id := utils.GetPaginationEmbedId(customId)
			userId := utils.GetPaginationEmbedUserId(id)

			if i.Member.User.ID != userId {
				return false
			}

			if utils.GetPaginationEmbed(id) == nil {
				return false
			}
		} else {
			return false
		}
		return true
	},
	Run: func(ctx *commands.ComponentContext) {
		customId := ctx.Inter.MessageComponentData().CustomID
		id := utils.GetPaginationEmbedId(customId)
		p := utils.GetPaginationEmbed(id)

		if strings.HasPrefix(customId, utils.PaginationEmbedPrev) {
			p.Prev(ctx.Inter)
		} else {
			p.Next(ctx.Inter)
		}
	},
}
