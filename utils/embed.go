package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

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

func GetDeclineContainer(components ...discordgo.MessageComponent) *discordgo.Container {
	c := &discordgo.Container{
		Components: []discordgo.MessageComponent{
			discordgo.TextDisplay{
				Content: "### ❌ 거부",
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

func GetUserIsNotRegisteredErrContainer(prefix string) *discordgo.Container {
	return GetErrorContainer(discordgo.TextDisplay{
		Content: fmt.Sprintf("해당 기능은 등록된 사용자만 쓸 수 있어요. `%s가입`으로 가입해주새요.", prefix),
	})
}
