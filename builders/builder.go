package builders

import "github.com/bwmarrin/discordgo"

type Builder[T any] interface {
	Build() discordgo.MessageComponent
}

type ComponentBuilder[T any] interface {
	Builder[T]
	AddComponents(components ...discordgo.MessageComponent) T
}
