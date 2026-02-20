package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders/customid"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(middlewares.CheckPaginationContainer())

		r.Component(customid.PaginationContainerFirst+"/{id}", func(e *handler.ComponentEvent) error {
			return builders.GetPaginationContainer(e.Vars["id"]).First(e)
		})

		r.Component(customid.PaginationContainerPrev+"/{id}", func(e *handler.ComponentEvent) error {
			return builders.GetPaginationContainer(e.Vars["id"]).Prev(e)
		})

		r.Component(customid.PaginationContainerNext+"/{id}", func(e *handler.ComponentEvent) error {
			return builders.GetPaginationContainer(e.Vars["id"]).Next(e)
		})

		r.Component(customid.PaginationContainerLast+"/{id}", func(e *handler.ComponentEvent) error {
			return builders.GetPaginationContainer(e.Vars["id"]).Last(e)
		})

		r.Component(customid.PaginationContainerPages+"/{id}", func(e *handler.ComponentEvent) error {
			return builders.GetPaginationContainer(e.Vars["id"]).ShowModal(e)
		})
	})
}
