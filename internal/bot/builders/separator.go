package builders

import "github.com/bwmarrin/discordgo"

type Separator struct {
	*discordgo.Separator
}

func SeparatorBuilder() *Separator {
	return &Separator{
		Separator: &discordgo.Separator{},
	}
}

func (s *Separator) SetDivider(divider bool) *Separator {
	s.Separator.Divider = &divider
	return s
}

func (s *Separator) SetSpacing(spacing discordgo.SeparatorSpacingSize) *Separator {
	s.Separator.Spacing = &spacing
	return s
}

func (s *Separator) Build() discordgo.MessageComponent {
	return s.Separator
}
