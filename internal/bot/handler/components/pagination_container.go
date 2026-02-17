package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(middlewares.CheckPaginationContainerMiddleware())

		r.Component(utils.PaginationContainerFirst+"/{id}", func(e *handler.ComponentEvent) error {
			return builders.GetPaginationContainer(e.Vars["id"]).First(e)
		})

		r.Component(utils.PaginationContainerPrev+"/{id}", func(e *handler.ComponentEvent) error {
			return builders.GetPaginationContainer(e.Vars["id"]).Prev(e)
		})

		r.Component(utils.PaginationContainerNext+"/{id}", func(e *handler.ComponentEvent) error {
			return builders.GetPaginationContainer(e.Vars["id"]).Next(e)
		})

		r.Component(utils.PaginationContainerLast+"/{id}", func(e *handler.ComponentEvent) error {
			return builders.GetPaginationContainer(e.Vars["id"]).Last(e)
		})

		r.Component(utils.PaginationContainerPages+"/{id}", func(e *handler.ComponentEvent) error {
			return builders.GetPaginationContainer(e.Vars["id"]).ShowModal(e)
		})
	})
}
