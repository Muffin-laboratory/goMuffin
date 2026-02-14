package loader

import "github.com/disgoorg/disgo/handler"

type Component struct {
	Middlewares handler.Middlewares
	Handle      func(r handler.Router)
}

func (d *Discommand) LoadComponent(c *Component) {
	d.Router.Group(func(r handler.Router) {
		r.Use(c.Middlewares...)
		c.Handle(r)
	})
}
