package loader

import (
	"context"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/bwmarrin/discordgo"
)

type Component struct {
	Parse             parse
	Run               run
	DeferredReply     bool
	DeferReplyOptions *discordgo.InteractionResponseData
	DeferredUpdate    bool
}

func (d *Discommand) LoadComponent(c *Component) {
	defer componentMutex.Unlock()
	componentMutex.Lock()
	d.Components = append(d.Components, c)
}

func (d *Discommand) ComponentRun(s *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var err error

	i := &builders.InteractionCreate{
		InteractionCreate: inter,
		Session:           s,
	}

	i.InteractionCreate.User = builders.GetInteractionUser(inter)

	for _, c := range d.Components {
		var ctx context.Context
		var cancel context.CancelFunc
		if c.DeferredReply || c.DeferredUpdate {
			ctx, cancel = context.WithTimeout(context.Background(), 15*time.Minute)
		} else {
			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		}
		i.Ctx = ctx

		defer cancel()

		if !c.Parse(i) {
			continue
		}

		if c.DeferredReply {
			if err := i.DeferReply(c.DeferReplyOptions); err != nil {
				return err
			}
		} else if c.DeferredUpdate {
			if err := i.DeferUpdate(); err != nil {
				return err
			}
		}

		err = c.Run(i)
		break
	}
	return err
}
