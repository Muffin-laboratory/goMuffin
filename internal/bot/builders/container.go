package builders

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Container struct {
	*discordgo.Container
}

func ContainerBuilder() *Container {
	return &Container{
		Container: &discordgo.Container{},
	}
}

func (c *Container) SetAccentColor(color int) *Container {
	c.Container.AccentColor = &color
	return c
}

func (c *Container) SetSpoiler(spoiler bool) *Container {
	c.Container.Spoiler = spoiler
	return c
}

func (c *Container) AddComponents(components ...ComponentBuilder) *Container {
	for _, cmp := range components {
		c.Container.Components = append(c.Container.Components, cmp.Build())
	}
	return c
}

func (c *Container) AddText(text string) *Container {
	c.AddComponents(TextDisplayBuilder(text))
	return c
}

func (c *Container) Build() discordgo.MessageComponent {
	return c.Container
}

func MakeErrorContainer(text string) *Container {
	return ContainerBuilder().
		AddComponents(
			TextDisplayBuilder("### ❌ 오류"),
			TextDisplayBuilder(text),
		)
}

func MakeDeclineContainer(text string) *Container {
	return ContainerBuilder().
		AddComponents(
			TextDisplayBuilder("### ❌ 거부"),
			TextDisplayBuilder(text),
		)
}

func MakeCanceledContainer(text string) *Container {
	return ContainerBuilder().
		AddComponents(
			TextDisplayBuilder("### ❌ 취소"),
			TextDisplayBuilder(text),
		)
}

func MakeSuccessContainer(text string) *Container {
	return ContainerBuilder().
		AddComponents(
			TextDisplayBuilder("### ✅ 성공"),
			TextDisplayBuilder(text),
		)
}

func MakeUserIsNotRegisteredErrContainer() *Container {
	return MakeErrorContainer("해당 기능은 등록된 사용자만 쓸 수 있어요. `/가입`으로 가입해주새요.")
}

func MakeUserIsBlockedContainer(globalName, reason string) *Container {
	return MakeDeclineContainer(fmt.Sprintf("- %s님은 서비스에서 차단되었어요.\n> 사유: %s", globalName, reason))
}
