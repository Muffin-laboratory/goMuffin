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
	desc    string
}

var PaginationEmbeds = make(map[string]*PaginationEmbed)

func makeComponents(id string, current, total int) *[]discordgo.MessageComponent {
	disabled := false

	if total == 1 {
		disabled = true
	}

	return &[]discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "이전",
					CustomID: MakePaginationEmbedPrev(id),
					Disabled: disabled,
				},
				discordgo.Button{
					Style:    discordgo.SecondaryButton,
					Label:    fmt.Sprintf("(%d/%d)", current, total),
					CustomID: MakePaginationEmbedPages(id),
					Disabled: disabled,
				},
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "다음",
					CustomID: MakePaginationEmbedNext(id),
					Disabled: disabled,
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
	}

	if len(data) <= 0 {
		p.Embed.Description = makeDesc(p.desc, "없음")
		p.Total = 1
	} else {
		p.Embed.Description = makeDesc(p.desc, data[0])
	}

	switch m := m.(type) {
	case *discordgo.MessageCreate:
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Reference:  m.Reference(),
			Embeds:     []*discordgo.MessageEmbed{p.Embed},
			Components: *makeComponents(id, p.Current, p.Total),
		})
	case *InteractionCreate:
		m.EditReply(&discordgo.WebhookEdit{
			Embeds:     &[]*discordgo.MessageEmbed{p.Embed},
			Components: makeComponents(id, p.Current, p.Total),
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

func (p *PaginationEmbed) Set(i *InteractionCreate, page int) {
	if page <= 0 {
		i.Reply(&discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "❌ 오류",
					Description: "해당 값은 0보다 커야해요.",
					Color:       EmbedFail,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	if page >= p.Total {
		i.Reply(&discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "❌ 오류",
					Description: "해당 값은 총 페이지의 수보다 작아야해요.",
					Color:       EmbedFail,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	p.Current = page

	p.Embed.Description = makeDesc(p.desc, p.Data[p.Current-1])

	i.Update(&discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{p.Embed},
		Components: *makeComponents(p.id, p.Current, p.Total),
	})
}

func (p *PaginationEmbed) ShowModal(i *InteractionCreate) {
	i.ShowModal(&ModalData{
		CustomId: MakePaginationEmbedModal(p.id),
		Title:    fmt.Sprintf("%s의 리스트", i.Session.State.User.Username),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    MakePaginationEmbedSetPage(p.id),
						Label:       "페이지",
						Style:       discordgo.TextInputShort,
						Placeholder: "이동할 페이지를 여기에 적어주세요.",
						Required:    true,
					},
				},
			},
		},
	})
}
