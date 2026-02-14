package loader

import "github.com/disgoorg/disgo/handler"

func (d *Discommand) RegisterHandler(handlerFunc func(r handler.Router)) {
	d.Router().Group(func(r handler.Router) {
		handlerFunc(r)
	})
}
