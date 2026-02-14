package builders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
)

// PaginationContainer is container with page
type PaginationContainer struct {
	Containers []discord.ContainerComponent
	Current    int
	Total      int
	ID         string
	m          creatableEvent
	timer      *time.Timer
	deferred   bool
}

type updatableEvent interface {
	discord.Interaction
	UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error
	CreateFollowupMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error)
}

type creatableEvent interface {
	discord.Interaction
	Client() *bot.Client
	CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error
}

var paginationContainers = make(map[string]*PaginationContainer)

const endDuration = 10 * time.Minute

// PaginationContainerBuilder creates a new PaginationContainer
func PaginationContainerBuilder(m creatableEvent, deferred bool) *PaginationContainer {
	userID := m.User().ID.String()

	id := fmt.Sprintf("%s:%d", userID, rand.Intn(100))
	return &PaginationContainer{
		Current:  1,
		ID:       id,
		m:        m,
		timer:    time.NewTimer(endDuration),
		deferred: deferred,
	}
}

func (p *PaginationContainer) waitTimerEnd() {
	<-p.timer.C
	delete(paginationContainers, p.ID)
}

func (p *PaginationContainer) resetTimer() {
	p.timer.Reset(endDuration)
}

// AddContainers adds containers
func (p *PaginationContainer) AddContainers(containers ...discord.ContainerComponent) *PaginationContainer {
	p.Total += len(containers)
	p.Containers = append(p.Containers, containers...)
	return p
}

// Start starts the paginated-container
func (p *PaginationContainer) Start() error {
	container := p.Containers[0].AddComponents(makeComponents(p.ID, p.Current, p.Total))
	paginationContainers[p.ID] = p

	go p.waitTimerEnd()

	if p.deferred {
		_, err := p.m.Client().Rest.UpdateInteractionResponse(
			p.m.ApplicationID(),
			p.m.Token(),
			discord.NewMessageUpdateV2([]discord.LayoutComponent{container}),
		)
		return err
	}

	return p.m.CreateMessage(
		discord.NewMessageCreateV2(container).
			WithEphemeral(true),
	)
}

func makeComponents(id string, current, total int) discord.ActionRowComponent {
	disabled := false

	if total == 1 {
		disabled = true
	}

	return discord.NewActionRow(
		discord.NewPrimaryButton("처음", utils.MakePaginationContainerFirst(id)).
			WithDisabled(disabled),
		discord.NewPrimaryButton("이전", utils.MakePaginationContainerPrev(id)).
			WithDisabled(disabled),
		discord.NewSecondaryButton(fmt.Sprintf("(%d/%d)", current, total), utils.MakePaginationContainerPages(id)).
			WithDisabled(disabled),
		discord.NewPrimaryButton("다음", utils.MakePaginationContainerNext(id)).
			WithDisabled(disabled),
		discord.NewPrimaryButton("마지막", utils.MakePaginationContainerLast(id)).
			WithDisabled(disabled),
	)
}

// GetPaginationContainer gets PaginationContainer
func GetPaginationContainer(id string) *PaginationContainer {
	if p, ok := paginationContainers[id]; ok {
		return p
	}
	return nil
}

// First moves to first page
func (p *PaginationContainer) First(i *handler.ComponentEvent) error {
	return p.Set(i, 1)
}

// Prev move to previous page
func (p *PaginationContainer) Prev(i *handler.ComponentEvent) error {
	if p.Current == 1 {
		p.Current = p.Total
	} else {
		p.Current -= 1
	}

	return p.Set(i, p.Current)
}

// Next moves to next page
func (p *PaginationContainer) Next(i *handler.ComponentEvent) error {
	if p.Current >= p.Total {
		p.Current = 1
	} else {
		p.Current += 1
	}

	return p.Set(i, p.Current)
}

// Last moves to last page
func (p *PaginationContainer) Last(i *handler.ComponentEvent) error {
	return p.Set(i, p.Total)
}

// Set sets to page
func (p *PaginationContainer) Set(i updatableEvent, page int) error {
	p.resetTimer()

	if page <= 0 {
		p.Current = 1
	} else if page > p.Total {
		p.Current = p.Total
	} else {
		p.Current = page
	}

	container := p.Containers[p.Current-1].AddComponents(makeComponents(p.ID, p.Current, p.Total))
	return i.UpdateMessage(
		discord.NewMessageUpdateV2([]discord.LayoutComponent{container}),
	)
}

// ShowModal show discord's modal
func (p *PaginationContainer) ShowModal(i *handler.ComponentEvent) error {
	return i.Modal(
		discord.NewModalCreate(
			utils.MakePaginationContainerModal(p.ID),
			"페이지 설정",
			[]discord.LayoutComponent{
				discord.NewLabel(
					"이동할 페이지",
					discord.NewShortTextInput(utils.PaginationContainerSetPage).
						WithPlaceholder("페이지 번호를 여기에 입력...").
						WithValue(fmt.Sprint(p.Current)),
				).
					WithDescription("이동할 페이지의 번호를 입력해 주세요."),
			},
		),
	)
}
