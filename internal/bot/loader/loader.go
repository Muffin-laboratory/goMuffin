package loader

import (
	"sync"

	"github.com/disgoorg/disgo/handler"
)

type Discommand struct {
	Router     handler.Router
	Commands   map[string]*Command
	Components []*Component
	Modals     []*Modal
}

var (
	commandMutex   sync.Mutex
	componentMutex sync.Mutex
	modalMutex     sync.Mutex
)

var once sync.Once

var instance *Discommand

func GetDiscommand() *Discommand {
	once.Do(func() {
		instance = &Discommand{
			Router:     handler.New(),
			Commands:   map[string]*Command{},
			Components: []*Component{},
			Modals:     []*Modal{},
		}
	})

	return instance
}
