package builders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

const EndDuration = time.Minute * 10

// A PaginationContainer is container with page
type PaginationContainer struct {
	Containers []*discordgo.Container
	timer      *time.Timer
	Current    int
	Total      int
	ID         string
	m          any
}

var paginationContainers = make(map[string]*PaginationContainer)

func PaginationContainerBuilder(m any) *PaginationContainer {
	var userID string
	var p PaginationContainer

	switch m := m.(type) {
	case *MessageCreate:
		userID = m.Author.ID
	case *InteractionCreate:
		userID = m.Member.User.ID
	}

	id := fmt.Sprintf("%s/%d", userID, rand.Intn(100))
	p = PaginationContainer{
		Current: 1,
		ID:      id,
		m:       m,
		timer:   time.NewTimer(EndDuration),
	}

	p.startTimer()

	return &p
}

func (p *PaginationContainer) startTimer() {
	go func() {
		<-p.timer.C
		delete(paginationContainers, p.ID)
	}()
}

func (p *PaginationContainer) resetTimer() {
	p.timer.Reset(EndDuration)
}

func (p *PaginationContainer) AddContainers(containers ...*Container) *PaginationContainer {
	p.Total += len(containers)
	for _, container := range containers {
		p.Containers = append(p.Containers, container.Build().(*discordgo.Container))
	}
	return p
}

func (p *PaginationContainer) Start() error {
	container := *p.Containers[0]
	container.Components = append(container.Components, makeComponents(p.ID, p.Current, p.Total))

	paginationContainers[p.ID] = p

	return NewMessageSender(p.m).
		AddComponents(container).
		SetReply(true).
		SetEphemeral(true).
		SetComponentsV2(true).
		Send()
}

func makeComponents(id string, current, total int) *discordgo.ActionsRow {
	disabled := false

	if total == 1 {
		disabled = true
	}

	return ActionsRowBuilder(
		ButtonBuilder().
			SetStyle(discordgo.PrimaryButton).
			SetLabel("처음").
			SetCustomID(utils.MakePaginationContainerFirst(id)).
			SetDisabled(disabled),
		ButtonBuilder().
			SetStyle(discordgo.PrimaryButton).
			SetLabel("이전").
			SetCustomID(utils.MakePaginationContainerPrev(id)).
			SetDisabled(disabled),
		ButtonBuilder().
			SetStyle(discordgo.SecondaryButton).
			SetLabel(fmt.Sprintf("(%d/%d)", current, total)).
			SetCustomID(utils.MakePaginationContainerPages(id)).
			SetDisabled(disabled),
		ButtonBuilder().
			SetStyle(discordgo.PrimaryButton).
			SetLabel("다음").
			SetCustomID(utils.MakePaginationContainerNext(id)).
			SetDisabled(disabled),
		ButtonBuilder().
			SetStyle(discordgo.PrimaryButton).
			SetLabel("마지막").
			SetCustomID(utils.MakePaginationContainerLast(id)).
			SetDisabled(disabled),
	).
		Build().(*discordgo.ActionsRow)
}

func GetPaginationContainer(id string) *PaginationContainer {
	if p, ok := paginationContainers[id]; ok {
		return p
	}
	return nil
}

func (p *PaginationContainer) First(i *InteractionCreate) error {
	return p.Set(i, 1)
}

func (p *PaginationContainer) Prev(i *InteractionCreate) error {
	if p.Current == 1 {
		p.Current = p.Total
	} else {
		p.Current -= 1
	}

	return p.Set(i, p.Current)
}

func (p *PaginationContainer) Next(i *InteractionCreate) error {
	if p.Current >= p.Total {
		p.Current = 1
	} else {
		p.Current += 1
	}

	return p.Set(i, p.Current)
}

func (p *PaginationContainer) Last(i *InteractionCreate) error {
	return p.Set(i, p.Total)
}

func (p *PaginationContainer) Set(i *InteractionCreate, page int) error {
	p.resetTimer()

	if page <= 0 {
		p.Current = 1
	} else if page > p.Total {
		p.Current = p.Total
	} else {
		p.Current = page
	}

	container := *p.Containers[p.Current-1]
	container.Components = append(container.Components, makeComponents(p.ID, p.Current, p.Total))

	return i.Update(&discordgo.InteractionResponseData{
		Flags:      discordgo.MessageFlagsIsComponentsV2,
		Components: []discordgo.MessageComponent{container},
	})
}

func (p *PaginationContainer) ShowModal(i *InteractionCreate) error {
	return i.ShowModal(&ModalData{
		CustomId: utils.MakePaginationContainerModal(p.ID),
		Title:    fmt.Sprintf("%s의 리스트", i.Session.State.User.Username),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    utils.MakePaginationContainerSetPage(p.ID),
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
