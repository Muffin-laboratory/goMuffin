package builders

import "github.com/bwmarrin/discordgo"

// A Separator is builder for separator
type Separator struct {
	separator *discordgo.Separator
}

// SeparatorBuilder create a new Separator
func SeparatorBuilder() *Separator {
	return &Separator{
		separator: &discordgo.Separator{},
	}
}

// SetDivider sets divider
func (s *Separator) SetDivider(divider bool) *Separator {
	s.separator.Divider = &divider
	return s
}

// SetSpacing sets spacing
func (s *Separator) SetSpacing(spacing discordgo.SeparatorSpacingSize) *Separator {
	s.separator.Spacing = &spacing
	return s
}

// Build returns discordgo.Separator(discordgo.MessageComponent)
func (s *Separator) Build() discordgo.MessageComponent {
	return s.separator
}
