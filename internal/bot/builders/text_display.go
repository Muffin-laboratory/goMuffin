package builders

import "github.com/bwmarrin/discordgo"

type TextDisplay struct {
	*discordgo.TextDisplay
}

func TextDisplayBuilder(text string) *TextDisplay {
	return &TextDisplay{
		TextDisplay: &discordgo.TextDisplay{
			Content: text,
		},
	}
}

func (t *TextDisplay) SetText(text string) *TextDisplay {
	t.TextDisplay.Content = text
	return t
}

func (t *TextDisplay) Build() discordgo.MessageComponent {
	return t.TextDisplay
}
