package modals

import (
	"strconv"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var PaginationContainerModal = &loader.Modal{
	Handle: func(r handler.Router) {
		r.Modal(utils.PaginationContainerModal+"/{id}", func(inter *handler.ModalEvent) error {
			p := builders.GetPaginationContainer(inter.Vars["id"])

			page, err := strconv.Atoi(inter.Data.Text(utils.PaginationContainerSetPage))
			if err != nil {
				return inter.CreateMessage(
					discord.NewMessageCreateBuilder().
						SetComponents(builders.MakeErrorContainer("해당 값은 숫자여야해요.")).
						SetIsComponentsV2(true).
						SetEphemeral(true).
						Build(),
				)
			}

			return p.Set(inter, page)
		})
	},
}

func init() {
	loader.GetDiscommand().LoadModal(PaginationContainerModal)
}
