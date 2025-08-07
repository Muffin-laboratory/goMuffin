package utils

import "regexp"

var (
	RegexpFlexibleString    = regexp.MustCompile(`[^\s"'「」«»]+|"([^"]*)"|'([^']*)'|「([^」]*)」|«([^»]*)»`)
	RegexpDecimals          = regexp.MustCompile(`\d+`)
	RegexpItemId            = regexp.MustCompile(`no=\d+`)
	RegexpUserId            = regexp.MustCompile(`user_id=\d+`)
	RegexpId                = regexp.MustCompile(`id=[^&]*`)
	RegexpDiscordEmoji      = regexp.MustCompile(`<a?:\w+:\d+>`)
	RegexpLearnQueryLength  = regexp.MustCompile(`개수:(\d+)`)
	RegexpPaginationEmbedId = regexp.MustCompile(`^(\d+)/(\d+)$`)
)
