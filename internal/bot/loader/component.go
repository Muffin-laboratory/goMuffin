package loader

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/events"
)

type componentFunc[T any] func(ctx context.Context, i *events.ComponentInteractionCreate) T

type Component struct {
	Parse            componentFunc[bool]
	Run              componentFunc[error]
	DeferredReply    bool
	IsDeferEphemeral bool
	DeferredUpdate   bool
}

func (d *Discommand) LoadComponent(c *Component) {
	defer componentMutex.Unlock()
	componentMutex.Lock()
	d.Components = append(d.Components, c)
}

func (d *Discommand) ComponentRun(i *events.ComponentInteractionCreate) error {
	var err error

	for _, c := range d.Components {
		var ctx context.Context
		var cancel context.CancelFunc
		if c.DeferredReply || c.DeferredUpdate {
			ctx, cancel = context.WithTimeout(context.Background(), 15*time.Minute)
		} else {
			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		}
		defer cancel()

		if !c.Parse(ctx, i) {
			continue
		}

		if c.DeferredReply {
			if err := i.DeferCreateMessage(c.IsDeferEphemeral); err != nil {
				return err
			}
		} else if c.DeferredUpdate {
			if err := i.DeferUpdateMessage(); err != nil {
				return err
			}
		}

		err = c.Run(ctx, i)
		break
	}
	return err
}
