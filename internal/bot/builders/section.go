package builders

import (
	"github.com/bwmarrin/discordgo"
)

// A Section is builder for section
type Section struct {
	section *discordgo.Section
}

// SectionBuilder creates a new section
func SectionBuilder() *Section {
	return &Section{
		section: &discordgo.Section{},
	}
}

// SetAccessory sets accessory
func (s *Section) SetAccessory(accessory ComponentBuilder) *Section {
	s.section.Accessory = accessory.Build()
	return s
}

// AddText adds a text
func (s *Section) AddText(format string, a ...any) *Section {
	s.section.Components = append(s.section.Components, TextDisplayBuilder(format, a...).Build())
	return s
}

// Build returns discordgo.Section(discordgo.MessageComponent)
func (s *Section) Build() discordgo.MessageComponent {
	return s.section
}
