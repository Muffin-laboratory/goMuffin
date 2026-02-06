package builders

import (
	"github.com/bwmarrin/discordgo"
)

// A Modal is a window to collect user's answer.
type Modal struct {
	CustomID   string                       `json:"custom_id"`
	Title      string                       `json:"title"`
	Components []discordgo.MessageComponent `json:"components"`
}

// ModalBuilder creates a new modal.
func ModalBuilder() *Modal {
	return &Modal{}
}

// SetCustomID sets its custom id.
func (m *Modal) SetCustomID(customID string) *Modal {
	m.CustomID = customID
	return m
}

// SetTitle sets its title.
func (m *Modal) SetTitle(title string) *Modal {
	m.Title = title
	return m
}

// AddComponents adds its components.
// Components should be a ActionsRow or Label component.
func (m *Modal) AddComponents(components ...ComponentBuilder) *Modal {
	for _, cmp := range components {
		m.Components = append(m.Components, cmp.Build())
	}

	return m
}
