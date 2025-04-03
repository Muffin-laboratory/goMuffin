package utils

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

func CodeBlockWithLanguage(language string, content string) string {
	return fmt.Sprintf("```%s\n%s\n```", language, content)
}

func CodeBlock(content string) string {
	return fmt.Sprintf("```\n%s\n```", content)
}

func Time(time *time.Time) string {
	return fmt.Sprintf("<t:%d>", time.Unix())
}

func TimeWithStyle(time *time.Time, style string) string {
	return fmt.Sprintf("<t:%d:%s>", time.Unix(), style)
}
