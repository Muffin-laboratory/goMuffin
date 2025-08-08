package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
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

func GetCanceledContainer(components ...discordgo.MessageComponent) *discordgo.Container {
	c := &discordgo.Container{
		Components: []discordgo.MessageComponent{
			discordgo.TextDisplay{
				Content: "### ❌ 취소",
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

func GetUserIsBlockedContainer(globalName, reason string) *discordgo.Container {
	return GetDeclineContainer(discordgo.TextDisplay{
		Content: fmt.Sprintf("- %s님은 서비스에서 차단되었어요.\n> 사유: %s", globalName, reason),
	})
}
