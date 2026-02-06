package builders

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// A TextDisplay is builder for text display
type TextDisplay struct {
	textDisplay *discordgo.TextDisplay
}

// TextDisplayBuilder creates a new TextDisplay
func TextDisplayBuilder(format string, a ...any) *TextDisplay {
	return &TextDisplay{
		textDisplay: &discordgo.TextDisplay{
			Content: fmt.Sprintf(format, a...),
		},
	}
}

// SetText sets text
func (t *TextDisplay) SetText(format string, a ...any) *TextDisplay {
	t.textDisplay.Content = fmt.Sprintf(format, a...)
	return t
}

// Build returns discordgo.TextDisplay(discordgo.MessageComponent)
func (t *TextDisplay) Build() discordgo.MessageComponent {
	return t.textDisplay
}
