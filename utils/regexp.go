package utils

import "regexp"

var (
	RegexpFlexibleString    = regexp.MustCompile(`[^\s"'「」«»]+|"([^"]*)"|'([^']*)'|「([^」]*)」|«([^»]*)»`)
	RegexpDecimals          = regexp.MustCompile(`\d+`)
	RegexpDLDItemId         = regexp.MustCompile(`no=\d+`)
	RegexpDLDUserId         = regexp.MustCompile(`user_id=\d+`)
	RegexpDLDId             = regexp.MustCompile(`id=[^&]*`)
	RegexpEmoji             = regexp.MustCompile(`<a?:\w+:\d+>`)
	RegexpLearnQueryCommand = regexp.MustCompile(`단어:([^\n대답개수:]*)`)
	RegexpLearnQueryResult  = regexp.MustCompile(`대답:([^\n단어개수:]*)`)
	RegexpLearnQueryLength  = regexp.MustCompile(`개수:(\d+)`)
	RegexpPaginationEmbedId = regexp.MustCompile(`^(\d+)/(\d+)$`)
)
