package loader

import (
	"sync"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type Discommand struct {
	router      handler.Router
	OldCommands map[string]*Command
	commands    []discord.ApplicationCommandCreate
	devCommands []discord.ApplicationCommandCreate
}

var (
	commandMutex sync.Mutex
)

var once sync.Once
var instance *Discommand
var timeout = 1 * time.Minute

func Timeout() time.Duration {
	return timeout
}

func GetDiscommand() *Discommand {
	once.Do(func() {
		r := handler.New()
		instance = &Discommand{
			router:      r,
			OldCommands: make(map[string]*Command),
		}

	})

	return instance
}

func (d *Discommand) Router() handler.Router {
	return d.router
}
