package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/handler"
)

var PaginationContainerComponent *loader.Component = &loader.Component{
	Handle: func(r handler.Router) {
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
	},
}

func init() {
	loader.GetDiscommand().LoadComponent(PaginationContainerComponent)
}
