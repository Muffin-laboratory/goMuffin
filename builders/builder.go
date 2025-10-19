package builders

import "github.com/bwmarrin/discordgo"

type ComponentBuilder interface {
	Build() discordgo.MessageComponent
}
