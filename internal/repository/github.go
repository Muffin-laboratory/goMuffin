package repository

import (
	"github.com/google/go-github/v74/github"
)

var ghInstance *github.Client

func GetGHClient() *github.Client {
	if ghInstance == nil {
		ghInstance = github.NewClient(nil)
	}

	return ghInstance
}
