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
	Id         string
	m          any
}

var PaginationEmbeds = make(map[string]*PaginationEmbed)

func PaginationEmbedBuilder(m any) *PaginationEmbed {
	var userId string

	switch m := m.(type) {
	case *MessageCreate:
		userId = m.Author.ID
	case *InteractionCreate:
		userId = m.Member.User.ID
	}

	id := fmt.Sprintf("%s/%d", userId, rand.Intn(100))
	return &PaginationEmbed{
		Current: 1,
		Id:      id,
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
	container.Components = append(container.Components, makeComponents(p.Id, p.Current, p.Total))

	PaginationEmbeds[p.Id] = p

	err := NewMessageSender(p.m).
		AddComponents(container).
		SetReply(true).
		SetEphemeral(true).
		SetComponentsV2(true).
		Send()
	return err
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
			Components: []discordgo.MessageComponent{
				GetErrorContainer(discordgo.TextDisplay{Content: "해당 페이지가 처음ㅇ이에요."}),
			},
			Flags: discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsIsComponentsV2,
		})
		return
	}

	p.Current -= 1

	p.Set(i, p.Current)
}

func (p *PaginationEmbed) Next(i *InteractionCreate) {
	if p.Current >= p.Total {
		i.Reply(&discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				GetErrorContainer(discordgo.TextDisplay{Content: "해당 페이지가 마지막ㅇ이에요."}),
			},
			Flags: discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsIsComponentsV2,
		})
		return
	}

	p.Current += 1

	p.Set(i, p.Current)
}

func (p *PaginationEmbed) Set(i *InteractionCreate, page int) error {
	if page <= 0 {
		i.Reply(&discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				GetErrorContainer(discordgo.TextDisplay{Content: "해당 값은 0보다 커야해요."}),
			},
			Flags: discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsIsComponentsV2,
		})
		return nil
	}

	if page > p.Total {
		i.Reply(&discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				GetErrorContainer(discordgo.TextDisplay{Content: "해당 값은 총 페이지의 수보다 작아야해요."}),
			},
			Flags: discordgo.MessageFlagsEphemeral | discordgo.MessageFlagsIsComponentsV2,
		})
		return nil
	}

	p.Current = page

	container := *p.Containers[p.Current-1]
	container.Components = append(container.Components, makeComponents(p.Id, p.Current, p.Total))

	err := i.Update(&discordgo.InteractionResponseData{
		Flags:      discordgo.MessageFlagsIsComponentsV2,
		Components: []discordgo.MessageComponent{container},
	})
	return err
}

func (p *PaginationEmbed) ShowModal(i *InteractionCreate) {
	i.ShowModal(&ModalData{
		CustomId: MakePaginationEmbedModal(p.Id),
		Title:    fmt.Sprintf("%s의 리스트", i.Session.State.User.Username),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    MakePaginationEmbedSetPage(p.Id),
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
