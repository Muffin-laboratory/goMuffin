package builders

import "github.com/bwmarrin/discordgo"

// A TextInput is a builder of text input.
type TextInput struct {
	textInput *discordgo.TextInput
}

// TextInputBuilder creates a new TextInput.
func TextInputBuilder() *TextInput {
	return &TextInput{
		textInput: &discordgo.TextInput{},
	}
}

// SetCustomID sets its custom id.
func (t *TextInput) SetCustomID(customID string) *TextInput {
	t.textInput.CustomID = customID
	return t
}

// SetStyle sets its style.
func (t *TextInput) SetStyle(style discordgo.TextInputStyle) *TextInput {
	t.textInput.Style = style
	return t
}

// SetMinLength sets its minimum length.
func (t *TextInput) SetMinLength(minLength int) *TextInput {
	t.textInput.MinLength = minLength
	return t
}

// SetMaxLength sets its maximum length.
func (t *TextInput) SetMaxLength(maxLength int) *TextInput {
	t.textInput.MaxLength = maxLength
	return t
}

// SetRequired sets the field is required.
func (t *TextInput) SetRequired(required bool) *TextInput {
	t.textInput.Required = required
	return t
}

// SetValue sets its pre-filled Value.
func (t *TextInput) SetValue(value string) *TextInput {
	t.textInput.Value = value
	return t
}

// SetPlaceholder sets its placeholder.
func (t *TextInput) SetPlaceholder(placeholder string) *TextInput {
	t.textInput.Placeholder = placeholder
	return t
}

// Deprecated: Use Label.SetLabel() and Label.SetDescription() instead.
// SetLabel sets its label.
func (t *TextInput) SetLabel(label string) *TextInput {
	t.textInput.Label = label
	return t
}

// Build returns a discordgo.TextInput(discordgo.MessageComponent).
func (t *TextInput) Build() discordgo.MessageComponent {
	return t.textInput
}
