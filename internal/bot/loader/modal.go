package loader

import "github.com/disgoorg/disgo/handler"

type Modal struct {
	Middlewares      handler.Middlewares
	Handle           func(r handler.Router)
	Deferred         bool
	IsDeferEphemeral bool
}

func (d *Discommand) LoadModal(m *Modal) {
	d.Router.Group(func(r handler.Router) {
		r.Use(m.Middlewares...)
		m.Handle(r)
	})
}
