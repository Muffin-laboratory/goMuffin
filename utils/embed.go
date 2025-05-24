package utils

import "github.com/bwmarrin/discordgo"

const (
	EmbedDefault int = 0xaddb87
	EmbedFail    int = 0xff0000
	EmbedSuccess int = 0x00ff00
)

func GetErrorContainer(components ...discordgo.MessageComponent) *discordgo.Container {
	c := &discordgo.Container{
		Components: []discordgo.MessageComponent{
			discordgo.TextDisplay{
				Content: "### ❌ 오류",
			},
		},
	}

	c.Components = append(c.Components, components...)
	return c
}

func GetSuccessContainer(components ...discordgo.MessageComponent) *discordgo.Container {
	c := &discordgo.Container{
		Components: []discordgo.MessageComponent{
			discordgo.TextDisplay{
				Content: "### ✅ 성공",
			},
		},
	}

	c.Components = append(c.Components, components...)
	return c
}
