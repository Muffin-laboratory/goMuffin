package loader

import "sync"

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
