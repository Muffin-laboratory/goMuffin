package builders

import "github.com/bwmarrin/discordgo"

// A Button is builder for button
type Button struct {
	button *discordgo.Button
}

// ButtonBuilder creates a new Button
func ButtonBuilder() *Button {
	return &Button{
		button: &discordgo.Button{},
	}
}

// SetLabel sets label
func (b *Button) SetLabel(label string) *Button {
	b.button.Label = label
	return b
}

// SetStyle sets button style
func (b *Button) SetStyle(style discordgo.ButtonStyle) *Button {
	b.button.Style = style
	return b
}

// SetDisabled sets disabled
func (b *Button) SetDisabled(disabled bool) *Button {
	b.button.Disabled = disabled
	return b
}

// SetEmoji sets emoji
func (b *Button) SetEmoji(emoji discordgo.ComponentEmoji) *Button {
	b.button.Emoji = &emoji
	return b
}

// SetURL sets url
func (b *Button) SetURL(url string) *Button {
	b.button.URL = url
	return b
}

// SetCustomID sets custom id
func (b *Button) SetCustomID(customID string) *Button {
	b.button.CustomID = customID
	return b
}

// SetSKUID sets sku id
func (b *Button) SetSKUID(skuID string) *Button {
	b.button.SKUID = skuID
	return b
}

// Build returns discordgo.Button(discordgo.MessageComponent)
func (b *Button) Build() discordgo.MessageComponent {
	return b.button
}
