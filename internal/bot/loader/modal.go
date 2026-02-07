package loader

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/events"
)

type modalFunc[T any] func(ctx context.Context, i *events.ModalSubmitInteractionCreate) T

type Modal struct {
	Parse            modalFunc[bool]
	Run              modalFunc[error]
	Deferred         bool
	IsDeferEphemeral bool
}

func (d *Discommand) LoadModal(m *Modal) {
	defer modalMutex.Unlock()
	modalMutex.Lock()
	d.Modals = append(d.Modals, m)
}

func (d *Discommand) ModalRun(i *events.ModalSubmitInteractionCreate) error {
	var err error

	for _, m := range d.Modals {
		var ctx context.Context
		var cancel context.CancelFunc
		if m.Deferred {
			ctx, cancel = context.WithTimeout(context.Background(), 15*time.Minute)
		} else {
			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		}
		defer cancel()

		if !m.Parse(ctx, i) {
			continue
		}

		if m.Deferred {
			if err := i.DeferCreateMessage(m.IsDeferEphemeral); err != nil {
				return err
			}
		}

		err = m.Run(ctx, i)
		break
	}
	return err
}
