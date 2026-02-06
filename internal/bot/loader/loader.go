package loader

import (
	"sync"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
)

type run func(inter *builders.InteractionCreate) error
type parse func(inter *builders.InteractionCreate) bool

type Discommand struct {
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
			Commands:   map[string]*Command{},
			Components: []*Component{},
			Modals:     []*Modal{},
		}
	})

	return instance
}
