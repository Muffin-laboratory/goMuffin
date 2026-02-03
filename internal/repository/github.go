package repository

import (
	"sync"

	"github.com/google/go-github/v74/github"
)

var ghInstance *github.Client
var ghOnce sync.Once

func GetGHClient() *github.Client {
	ghOnce.Do(func() {
		ghInstance = github.NewClient(nil)
	})

	return ghInstance
}
