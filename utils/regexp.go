package utils

import "regexp"

var FlexibleStringParser *regexp.Regexp = regexp.MustCompile("[^\\s\"'「」«»]+|\"([^\"]*)\"|'([^']*)'|「([^」]*)」|«([^»]*)»")
var Decimals *regexp.Regexp = regexp.MustCompile(`\d+`)
var ItemIdRegexp *regexp.Regexp = regexp.MustCompile(`No.\d+`)
var EmojiRegexp *regexp.Regexp = regexp.MustCompile(`<a?:\w+:\d+>`)
