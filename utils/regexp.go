package utils

import "regexp"

var (
	FlexibleStringParser = regexp.MustCompile("[^\\s\"'「」«»]+|\"([^\"]*)\"|'([^']*)'|「([^」]*)」|«([^»]*)»")
	Decimals             = regexp.MustCompile(`\d+`)
	ItemIdRegexp         = regexp.MustCompile(`No.\d+`)
	EmojiRegexp          = regexp.MustCompile(`<a?:\w+:\d+>`)
	LearnQueryCommand    = regexp.MustCompile(`^단어:`)
	LearnQueryResult     = regexp.MustCompile(`^대답:`)
)
