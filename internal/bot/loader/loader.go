package loader

import (
	"sync"
	"time"

	"github.com/disgoorg/disgo/handler"
)

type Discommand struct {
	Router   handler.Router
	Commands map[string]*Command
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
			Router:   r,
			Commands: make(map[string]*Command),
		}

	})

	return instance
}
