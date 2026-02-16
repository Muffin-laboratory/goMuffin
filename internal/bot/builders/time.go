package builders

import (
	"fmt"
	"time"
)

const (
	ShortTime = "t"
	LongTime  = "T"

	ShortDate = "d"
	LongDate  = "D"

	ShortDateTime = "f"
	LongDateTime  = "F"

	RelativeTime = "R"
)

func Time(time *time.Time, style string) string {
	if style == "" {
		return fmt.Sprintf("<t:%d>", time.Unix())
	}
	return fmt.Sprintf("<t:%d:%s>", time.Unix(), style)
}
