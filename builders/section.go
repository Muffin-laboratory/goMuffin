package builders

import "github.com/bwmarrin/discordgo"

type Section struct {
	*discordgo.Section
}

func SectionBuilder() *Section {
	return &Section{
		Section: &discordgo.Section{},
	}
}

func (s *Section) SetAccessory(accessory discordgo.MessageComponent) *Section {
	s.Section.Accessory = accessory
	return s
}

func (s *Section) AddComponents(components ...discordgo.MessageComponent) *Section {
	s.Section.Components = append(s.Section.Components, components...)
	return s
}

func (s *Section) Build() discordgo.MessageComponent {
	return s.Section
}
