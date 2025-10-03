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

func (s *Section) SetAccessory(accessory ComponentBuilder) *Section {
	s.Section.Accessory = accessory.Build()
	return s
}

func (s *Section) AddText(text string) *Section {
	s.Components = append(s.Components, discordgo.TextDisplay{
		Content: text,
	})
	return s
}

func (s *Section) Build() discordgo.MessageComponent {
	return s.Section
}
