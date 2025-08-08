package utils

import "regexp"

var (
	RegexpFlexibleString    = regexp.MustCompile(`[^\s"'「」«»]+|"([^"]*)"|'([^']*)'|「([^」]*)」|«([^»]*)»`)
	RegexpDecimals          = regexp.MustCompile(`\d+`)
	RegexpItemID            = regexp.MustCompile(`no=\d+`)
	RegexpUserID            = regexp.MustCompile(`user_id=\d+`)
	RegexpID                = regexp.MustCompile(`id=[^&]*`)
	RegexpDiscordEmoji      = regexp.MustCompile(`<a?:\w+:\d+>`)
	RegexpLearnQueryLength  = regexp.MustCompile(`개수:(\d+)`)
	RegexpPaginationEmbedID = regexp.MustCompile(`^(\d+)/(\d+)$`)
)
