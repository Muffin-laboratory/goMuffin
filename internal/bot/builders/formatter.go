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

func InlineCode(content string) string {
	return fmt.Sprintf("`%s`", content)
}

func CodeBlock(language string, content string) string {
	if content == "" {
		return fmt.Sprintf("```\n%s\n```", language)
	}
	return fmt.Sprintf("```%s\n%s\n```", language, content)
}

func Time(time *time.Time, style string) string {
	if style == "" {
		return fmt.Sprintf("<t:%d>", time.Unix())
	}
	return fmt.Sprintf("<t:%d:%s>", time.Unix(), style)
}
