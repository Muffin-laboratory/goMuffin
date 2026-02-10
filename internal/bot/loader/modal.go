package loader

import "github.com/disgoorg/disgo/handler"

type Modal struct {
	Handle           func(r handler.Router)
	Deferred         bool
	IsDeferEphemeral bool
}

func (d *Discommand) LoadModal(m *Modal) {
	defer modalMutex.Unlock()
	modalMutex.Lock()
	d.Modals = append(d.Modals, m)
	m.Handle(d.Router)
}
