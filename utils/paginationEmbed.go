package utils

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

// PaginationEmbed is embed with page
type PaginationEmbed struct {
	Embed   *discordgo.MessageEmbed
	Data    []string
	Current int
	Total   int
	id      string
	s       *discordgo.Session
	desc    string
}

var PaginationEmbeds = make(map[string]*PaginationEmbed)

func makeComponents(id string, current, total int) *[]discordgo.MessageComponent {

	return &[]discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "이전",
					CustomID: MakePaginationEmbedPrev(id),
					Disabled: false,
				},
				discordgo.Button{
					Style:    discordgo.SecondaryButton,
					Label:    fmt.Sprintf("(%d/%d)", current, total),
					CustomID: MakePaginationEmbedPages(id, current, total),
					Disabled: true,
				},
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "다음",
					CustomID: MakePaginationEmbedNext(id),
					Disabled: false,
				},
			},
		},
	}
}

func makeDesc(desc, item string) string {
	var newDesc string

	if desc == "" {
		newDesc = item
	} else {
		newDesc = fmt.Sprintf(desc, item)
	}
	return newDesc
}

// StartPaginationEmbed starts new PaginationEmbed struct
func StartPaginationEmbed(s *discordgo.Session, m any, e *discordgo.MessageEmbed, data []string, defaultDesc string) {
	var userId string

	switch m := m.(type) {
	case *discordgo.MessageCreate:
		userId = m.Author.ID
	case *InteractionCreate:
		userId = m.Member.User.ID
	}

	id := fmt.Sprintf("%s/%d", userId, rand.Intn(12))
	p := &PaginationEmbed{
		Embed:   e,
		Data:    data,
		Current: 1,
		Total:   len(data),
		id:      id,
		s:       s,
		desc:    defaultDesc,
	}

	p.Embed.Description = makeDesc(p.desc, data[0])

	switch m := m.(type) {
	case *discordgo.MessageCreate:
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Reference:  m.Reference(),
			Embeds:     []*discordgo.MessageEmbed{p.Embed},
			Components: *makeComponents(id, 1, len(data)),
		})
	case *InteractionCreate:
		m.Reply(&discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{p.Embed},
			Components: *makeComponents(id, 1, len(data)),
		})
	}

	PaginationEmbeds[id] = p
}

func GetPaginationEmbed(id string) *PaginationEmbed {
	if p, ok := PaginationEmbeds[id]; ok {
		return p
	}
	return nil
}

func (p *PaginationEmbed) Prev(i *InteractionCreate) {
	if p.Current == 1 {
		i.Reply(&discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "❌ 오류",
					Description: "해당 페이지가 처음ㅇ이에요.",
					Color:       EmbedFail,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	p.Current -= 1

	p.Embed.Description = makeDesc(p.desc, p.Data[p.Current-1])

	i.Update(&discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{p.Embed},
		Components: *makeComponents(p.id, p.Current, p.Total),
	})
}

func (p *PaginationEmbed) Next(i *InteractionCreate) {
	if p.Current >= p.Total {
		i.Reply(&discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "❌ 오류",
					Description: "해당 페이지가 마지막ㅇ이에요.",
					Color:       EmbedFail,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	p.Current += 1

	p.Embed.Description = makeDesc(p.desc, p.Data[p.Current-1])

	i.Update(&discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{p.Embed},
		Components: *makeComponents(p.id, p.Current, p.Total),
	})
}
