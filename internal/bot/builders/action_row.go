package builders

import "github.com/bwmarrin/discordgo"

type ActionsRow struct {
	*discordgo.ActionsRow
}

func ActionsRowBuilder(components ...ComponentBuilder) *ActionsRow {
	row := &ActionsRow{
		ActionsRow: &discordgo.ActionsRow{},
	}
	row.AddComponents(components...)

	return row
}

func (r *ActionsRow) AddComponents(components ...ComponentBuilder) *ActionsRow {
	for _, cmp := range components {
		r.ActionsRow.Components = append(r.ActionsRow.Components, cmp.Build())
	}

	return r
}

func (r *ActionsRow) Build() discordgo.MessageComponent {
	return r.ActionsRow
}
