package utils

import "github.com/bwmarrin/discordgo"

func GetStyleFromBool(k bool) discordgo.ButtonStyle {
	if k {
		return discordgo.SuccessButton
	}

	return discordgo.SecondaryButton
}

func BoolToString(k bool) string {
	if k {
		return "활성화"
	}

	return "비활성화"
}
