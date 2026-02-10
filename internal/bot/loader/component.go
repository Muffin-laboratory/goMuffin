package loader

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

type componentFunc[T any] func(ctx context.Context, i *events.ComponentInteractionCreate) T

type Component struct {
	Middlewares      handler.Middlewares
	Handle           func(r handler.Router)
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
	if c.Handle != nil {
		d.Router.Group(func(r handler.Router) {
			r.Use(c.Middlewares...)
			c.Handle(r)
		})
	}
}

func (d *Discommand) ComponentRun(i *events.ComponentInteractionCreate) error {
	var err error

	for _, c := range d.Components {
		if c.Handle != nil {
			continue
		}

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
