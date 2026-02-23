package builders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders/customid"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
)

// PaginatedContainer is container with page
type PaginatedContainer struct {
	containers []discord.ContainerComponent
	current    int
	total      int
	id         string
	timer      *time.Timer
	deferred   bool
	event      creatableEvent
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

var paginatedContainers = make(map[string]*PaginatedContainer)

const endDuration = 10 * time.Minute

// NewPaginatedContainer creates a new PaginationContainer
func NewPaginatedContainer(event creatableEvent, deferred bool) *PaginatedContainer {
	id := fmt.Sprintf("%d:%d", event.User().ID, rand.Intn(100))
	return &PaginatedContainer{
		current:  1,
		id:       id,
		timer:    time.NewTimer(endDuration),
		deferred: deferred,
		event:    event,
	}
}

func (p *PaginatedContainer) waitTimerEnd() {
	<-p.timer.C
	p.event.Client().Rest.UpdateInteractionResponse(
		p.event.ApplicationID(),
		p.event.Token(),
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			p.containers[p.current-1],
		}),
	)
	delete(paginatedContainers, p.id)
}

func (p *PaginatedContainer) resetTimer() {
	p.timer.Reset(endDuration)
}

func (p *PaginatedContainer) setPage(page int) {
	if page <= 0 {
		p.current = 1
	} else if page > p.total {
		p.current = p.total
	} else {
		p.current = page
	}
}

// AddContainers adds containers
func (p *PaginatedContainer) AddContainers(containers ...discord.ContainerComponent) *PaginatedContainer {
	p.total += len(containers)
	p.containers = append(p.containers, containers...)
	return p
}

// SetStartPage sets start page number.
func (p *PaginatedContainer) SetStartPage(startPage int) *PaginatedContainer {
	p.setPage(startPage)
	return p
}

// Start starts the paginated-container
func (p *PaginatedContainer) Start() error {
	if len(p.containers) == 0 {
		return nil
	}

	container := p.containers[p.current-1].AddComponents(p.makeComponents())
	paginatedContainers[p.id] = p

	go p.waitTimerEnd()

	if p.deferred {
		_, err := p.event.Client().Rest.UpdateInteractionResponse(
			p.event.ApplicationID(),
			p.event.Token(),
			discord.NewMessageUpdateV2([]discord.LayoutComponent{container}),
		)
		return err
	}

	return p.event.CreateMessage(
		discord.NewMessageCreateV2(container).
			WithEphemeral(true),
	)
}

func (p *PaginatedContainer) makeComponents() discord.ActionRowComponent {
	disabled := false

	if p.total == 1 {
		disabled = true
	}

	return discord.NewActionRow(
		discord.NewPrimaryButton("", customid.MakePaginationContainerFirst(p.id)).
			WithEmoji(discord.NewComponentEmoji("⏪")).
			WithDisabled(disabled),
		discord.NewPrimaryButton("", customid.MakePaginationContainerPrev(p.id)).
			WithEmoji(discord.NewComponentEmoji("◀️")).
			WithDisabled(disabled),
		discord.NewSecondaryButton(fmt.Sprintf("(%d/%d)", p.current, p.total), customid.MakePaginationContainerPages(p.id)).
			WithDisabled(disabled),
		discord.NewPrimaryButton("", customid.MakePaginationContainerNext(p.id)).
			WithEmoji(discord.NewComponentEmoji("▶️")).
			WithDisabled(disabled),
		discord.NewPrimaryButton("", customid.MakePaginationContainerLast(p.id)).
			WithEmoji(discord.NewComponentEmoji("⏩")).
			WithDisabled(disabled),
	)
}

// GetPaginationContainer gets PaginationContainer
func GetPaginationContainer(id string) *PaginatedContainer {
	if p, ok := paginatedContainers[id]; ok {
		return p
	}
	return nil
}

// First moves to first page
func (p *PaginatedContainer) First(i *handler.ComponentEvent) error {
	if p.current == 1 {
		return p.Set(i, p.total)
	}

	return p.Set(i, 1)
}

// Prev move to previous page
func (p *PaginatedContainer) Prev(i *handler.ComponentEvent) error {
	if p.current == 1 {
		return p.Set(i, p.total)
	}

	return p.Set(i, p.current-1)
}

// Next moves to next page
func (p *PaginatedContainer) Next(i *handler.ComponentEvent) error {
	if p.current == p.total {
		return p.Set(i, 1)
	}

	return p.Set(i, p.current+1)
}

// Last moves to last page
func (p *PaginatedContainer) Last(i *handler.ComponentEvent) error {
	if p.current == p.total {
		return p.Set(i, 1)
	}

	return p.Set(i, p.total)
}

// Set sets to page
func (p *PaginatedContainer) Set(i updatableEvent, page int) error {
	p.resetTimer()
	p.setPage(page)

	container := p.containers[p.current-1].AddComponents(p.makeComponents())
	return i.UpdateMessage(
		discord.NewMessageUpdateV2([]discord.LayoutComponent{container}),
	)
}

// ShowModal show discord's modal
func (p *PaginatedContainer) ShowModal(i *handler.ComponentEvent) error {
	return i.Modal(
		discord.NewModalCreate(
			customid.MakePaginationContainerModal(p.id),
			"페이지 설정",
			[]discord.LayoutComponent{
				discord.NewLabel(
					"이동할 페이지",
					discord.NewShortTextInput(customid.PaginationContainerSetPage).
						WithPlaceholder("페이지 번호를 여기에 입력...").
						WithValue(fmt.Sprint(p.current)),
				).
					WithDescription("이동할 페이지의 번호를 입력해 주세요."),
			},
		),
	)
}
