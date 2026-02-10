package loader

import (
	"context"
	"sync"
	"time"

	"github.com/disgoorg/disgo/handler"
)

type Discommand struct {
	Router     handler.Router
	Commands   map[string]*Command
	Components []*Component
}

var (
	commandMutex   sync.Mutex
	componentMutex sync.Mutex
)

var once sync.Once

var instance *Discommand

func GetDiscommand() *Discommand {
	once.Do(func() {
		r := handler.New()
		r.DefaultContext(func() context.Context {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			return ctx
		})
		instance = &Discommand{
			Router:   r,
			Commands: make(map[string]*Command),
		}

	})

	return instance
}
