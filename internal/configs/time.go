package configs

import "time"

var StartedAt *time.Time

func init() {
	now := time.Now()
	StartedAt = &now
}
