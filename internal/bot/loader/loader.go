package loader

import (
	"sync"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type Discommand struct {
	router      handler.Router
	commands    []discord.ApplicationCommandCreate
	devCommands []discord.ApplicationCommandCreate
}

var once sync.Once
var instance *Discommand
var timeout = 1 * time.Minute

func Timeout() time.Duration {
	return timeout
}

func GetDiscommand() *Discommand {
	once.Do(func() {
		r := handler.New()
		r.Use(middlewares.SendErrorMessage())
		instance = &Discommand{
			router: r,
		}

	})

	return instance
}

func (d *Discommand) Router() handler.Router {
	return d.router
}

func (d *Discommand) RegisterHandler(handlerFunc func(r handler.Router)) {
	d.Router().Group(func(r handler.Router) {
		handlerFunc(r)
	})
}
