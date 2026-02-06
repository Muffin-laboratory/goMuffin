package builders

import "github.com/bwmarrin/discordgo"

// A Label is a builder for label
type Label struct {
	label *discordgo.Label
}

// LabelBuilder creates a new label
func LabelBuilder() *Label {
	return &Label{
		label: &discordgo.Label{},
	}
}

// SetLabel sets a label
func (l *Label) SetLabel(label string) *Label {
	l.label.Label = label
	return l
}

// SetDescription sets a description
func (l *Label) SetDescription(description string) *Label {
	l.label.Description = description
	return l
}

// SetComponent sets a component
func (l *Label) SetComponent(components ComponentBuilder) *Label {
	l.label.Component = components.Build()
	return l
}

// Build returns a discordgo.Label(discordgo.MessageComponent)
func (l *Label) Build() discordgo.MessageComponent {
	return l.label
}
