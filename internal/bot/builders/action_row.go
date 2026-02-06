package builders

import "github.com/bwmarrin/discordgo"

// An ActionsRow is contain buttons or a select menu.
type ActionsRow struct {
	*discordgo.ActionsRow
}

// ActionsRowBuilder creates actions row
func ActionsRowBuilder(components ...ComponentBuilder) *ActionsRow {
	row := &ActionsRow{
		ActionsRow: &discordgo.ActionsRow{},
	}
	row.AddComponents(components...)

	return row
}

// AddComponents adds components
func (r *ActionsRow) AddComponents(components ...ComponentBuilder) *ActionsRow {
	for _, cmp := range components {
		r.ActionsRow.Components = append(r.ActionsRow.Components, cmp.Build())
	}

	return r
}

// Build returns discordgo.ActionsRow(discordgo.MessageComponent)
func (r *ActionsRow) Build() discordgo.MessageComponent {
	return r.ActionsRow
}
