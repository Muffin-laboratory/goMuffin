// builders contains codes for builder
package builders

import "github.com/bwmarrin/discordgo"

// A ComponentBuilder is builder for discord message component
type ComponentBuilder interface {
	// Build returns to discordgo.MessageComponent
	Build() discordgo.MessageComponent
}
