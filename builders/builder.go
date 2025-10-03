package builders

import "github.com/bwmarrin/discordgo"

type ComponentBuilder[T any] interface {
	Build() discordgo.MessageComponent
	AddComponents(components ...discordgo.MessageComponent) T
}
