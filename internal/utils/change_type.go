package utils

import (
	"github.com/disgoorg/disgo/discord"
)

func GetStyleFromBool(k bool) discord.ButtonStyle {
	if k {
		return discord.ButtonStyleSuccess
	}

	return discord.ButtonStyleSecondary
}

func BoolToString(k bool) string {
	if k {
		return "활성화"
	}

	return "비활성화"
}
