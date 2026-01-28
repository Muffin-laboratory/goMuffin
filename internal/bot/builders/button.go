package builders

import "github.com/bwmarrin/discordgo"

type Button struct {
	*discordgo.Button
}

func ButtonBuilder() *Button {
	return &Button{
		Button: &discordgo.Button{},
	}
}

func (b *Button) SetLabel(label string) *Button {
	b.Button.Label = label
	return b
}

func (b *Button) SetStyle(style discordgo.ButtonStyle) *Button {
	b.Button.Style = style
	return b
}

func (b *Button) SetDisabled(disabled bool) *Button {
	b.Button.Disabled = disabled
	return b
}

func (b *Button) SetEmoji(emoji discordgo.ComponentEmoji) *Button {
	b.Button.Emoji = &emoji
	return b
}

func (b *Button) SetURL(url string) *Button {
	b.Button.URL = url
	return b
}

func (b *Button) SetCustomID(customID string) *Button {
	b.Button.CustomID = customID
	return b
}

func (b *Button) SetSKUID(skuID string) *Button {
	b.Button.SKUID = skuID
	return b
}

func (b *Button) Build() discordgo.MessageComponent {
	return b.Button
}
