package loader

import (
	"context"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/bwmarrin/discordgo"
)

type Modal struct {
	Parse        parse
	Run          run
	Deferred     bool
	DeferOptions *discordgo.InteractionResponseData
}

func (d *Discommand) LoadModal(m *Modal) {
	defer modalMutex.Unlock()
	modalMutex.Lock()
	d.Modals = append(d.Modals, m)
}

func (d *Discommand) ModalRun(s *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var err error

	i := &builders.InteractionCreate{
		InteractionCreate: inter,
		Session:           s,
	}

	for _, m := range d.Modals {
		var ctx context.Context
		var cancel context.CancelFunc
		if m.Deferred {
			ctx, cancel = context.WithTimeout(context.Background(), 15*time.Minute)
		} else {
			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		}
		i.Ctx = ctx

		defer cancel()

		if !m.Parse(i) {
			continue
		}

		if m.Deferred {
			if err := i.DeferReply(m.DeferOptions); err != nil {
				return err
			}
		}

		err = m.Run(i)
		break
	}
	return err
}
