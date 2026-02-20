package modals

import (
	"strconv"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders/customid"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(middlewares.CheckPaginationContainer())

		r.Modal(customid.PaginationContainerModal+"/{id}", func(inter *handler.ModalEvent) error {
			p := builders.GetPaginationContainer(inter.Vars["id"])

			page, err := strconv.Atoi(inter.Data.Text(customid.PaginationContainerSetPage))
			if err != nil {
				return inter.CreateMessage(
					discord.NewMessageCreateV2(builders.MakeErrorContainer("해당 값은 숫자여야해요.")).
						WithEphemeral(true),
				)
			}

			return p.Set(inter, page)
		})
	})
}
