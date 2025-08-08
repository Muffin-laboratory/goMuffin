package utils

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

// PaginationEmbed is embed with page
type PaginationEmbed struct {
	Container  *discordgo.Container
	Containers []*discordgo.Container
	Current    int
	Total      int
	ID         string
	m          any
}

var paginationEmbeds = make(map[string]*PaginationEmbed)

func PaginationEmbedBuilder(m any) *PaginationEmbed {
	var userID string

	switch m := m.(type) {
	case *MessageCreate:
		userID = m.Author.ID
	case *InteractionCreate:
		userID = m.Member.User.ID
	}

	id := fmt.Sprintf("%s/%d", userID, rand.Intn(100))
	return &PaginationEmbed{
		Current: 1,
		ID:      id,
		m:       m,
	}
}

func (p *PaginationEmbed) SetContainer(container discordgo.Container) *PaginationEmbed {
	p.Container = &container
	return p
}

func (p *PaginationEmbed) AddContainers(container ...*discordgo.Container) *PaginationEmbed {
	p.Total += len(container)
	p.Containers = append(p.Containers, container...)
	return p
}

func (p *PaginationEmbed) Start() error {
	return startPaginationEmbed(p)
}

func makeComponents(id string, current, total int) *discordgo.ActionsRow {
	disabled := false

	if total == 1 {
		disabled = true
	}

	return &discordgo.ActionsRow{
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
	}
}

func MakeDesc(desc, item string) string {
	var newDesc string

	if desc == "" {
		newDesc = item
	} else {
		newDesc = fmt.Sprintf(desc, item)
	}
	return newDesc
}

func startPaginationEmbed(p *PaginationEmbed) error {
	container := *p.Containers[0]
	container.Components = append(container.Components, makeComponents(p.ID, p.Current, p.Total))

	paginationEmbeds[p.ID] = p

	err := NewMessageSender(p.m).
		AddComponents(container).
		SetReply(true).
		SetEphemeral(true).
		SetComponentsV2(true).
		Send()
	return err
}

func GetPaginationEmbed(id string) *PaginationEmbed {
	if p, ok := paginationEmbeds[id]; ok {
		return p
	}
	return nil
}

func (p *PaginationEmbed) Prev(i *InteractionCreate) error {
	if p.Current == 1 {
		p.Current = p.Total
	} else {
		p.Current -= 1
	}

	return p.Set(i, p.Current)
}

func (p *PaginationEmbed) Next(i *InteractionCreate) error {
	if p.Current >= p.Total {
		p.Current = 1
	} else {
		p.Current += 1
	}

	return p.Set(i, p.Current)
}

func (p *PaginationEmbed) Set(i *InteractionCreate, page int) error {
	if page <= 0 {
		p.Current = 1
	} else if page > p.Total {
		p.Current = p.Total
	} else {
		p.Current = page
	}

	container := *p.Containers[p.Current-1]
	container.Components = append(container.Components, makeComponents(p.ID, p.Current, p.Total))

	err := i.Update(&discordgo.InteractionResponseData{
		Flags:      discordgo.MessageFlagsIsComponentsV2,
		Components: []discordgo.MessageComponent{container},
	})
	return err
}

func (p *PaginationEmbed) ShowModal(i *InteractionCreate) {
	i.ShowModal(&ModalData{
		CustomId: MakePaginationEmbedModal(p.ID),
		Title:    fmt.Sprintf("%s의 리스트", i.Session.State.User.Username),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    MakePaginationEmbedSetPage(p.ID),
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
