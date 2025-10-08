package builders

import "github.com/bwmarrin/discordgo"

type ActionsRow struct {
	*discordgo.ActionsRow
}

func ActionsRowBuilder(components ...discordgo.MessageComponent) *ActionsRow {
	row := ActionsRow{}
	row.ActionsRow.Components = append(row.ActionsRow.Components, components...)

	return &row
}

func (r *ActionsRow) Build() discordgo.MessageComponent {
	return r
}
