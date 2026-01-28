package commands

import (
	"context"
	"sync"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/bwmarrin/discordgo"
)

type run func(inter *builders.InteractionCreate) error
type parse func(inter *builders.InteractionCreate) bool

type CommandFlags uint8

type Command struct {
	*discordgo.ApplicationCommand
	DeferOptions *discordgo.InteractionResponseData
	Flags        CommandFlags
	Deferred     bool
	Run          run
	Autocomplete run
}

type Discommand struct {
	Commands   map[string]*Command
	Components []*Component
	Modals     []*Modal
}

type Component struct {
	Parse             parse
	Run               run
	DeferredReply     bool
	DeferReplyOptions *discordgo.InteractionResponseData
	DeferredUpdate    bool
}

type Modal struct {
	Parse        parse
	Run          run
	Deferred     bool
	DeferOptions *discordgo.InteractionResponseData
}

const (
	CommandFlagsIsRegistered CommandFlags = 1 << iota
	CommandFlagsIsBlocked
	CommandFlagsIsDeveloperOnlyCommand
)

var (
	commandMutex   sync.Mutex
	componentMutex sync.Mutex
	modalMutex     sync.Mutex
)

var instance *Discommand

func GetDiscommand() *Discommand {
	if instance == nil {
		instance = &Discommand{
			Commands:   map[string]*Command{},
			Components: []*Component{},
			Modals:     []*Modal{},
		}
	}

	return instance
}

func (d *Discommand) LoadCommand(c *Command) {
	defer commandMutex.Unlock()
	commandMutex.Lock()
	d.Commands[c.Name] = c
}

func (d *Discommand) LoadComponent(c *Component) {
	defer componentMutex.Unlock()
	componentMutex.Lock()
	d.Components = append(d.Components, c)
}

func (d *Discommand) LoadModal(m *Modal) {
	defer modalMutex.Unlock()
	modalMutex.Lock()
	d.Modals = append(d.Modals, m)
}

func (d *Discommand) ChatInputRun(name string, s *discordgo.Session, inter *discordgo.InteractionCreate) error {
	i := &builders.InteractionCreate{
		InteractionCreate: inter,
		Session:           s,
		Options:           builders.MakeCommandInteractionOptionsMap(inter.ApplicationCommandData().Options),
	}

	i.InteractionCreate.User = builders.GetInteractionUser(inter)

	if command, ok := d.Commands[name]; ok {
		var ctx context.Context
		var cancel context.CancelFunc
		if command.Deferred {
			ctx, cancel = context.WithTimeout(context.Background(), 15*time.Minute)
		} else {
			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		}
		i.Ctx = ctx

		defer cancel()

		if command.Deferred {
			if err := i.DeferReply(command.DeferOptions); err != nil {
				return err
			}
		}

		if command.Flags&CommandFlagsIsDeveloperOnlyCommand != 0 && i.User.ID != configs.GetConfig().Bot.OwnerID {
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeDeclineContainer("이 명령어는 개발자 전용 명령어에요.")).
				SetComponentsV2(true).
				SetEphemeral(true).
				Send()
		}

		if command.Flags&CommandFlagsIsRegistered != 0 && !repository.GetDatabase().Users.IsUser(ctx, i.User.ID) {
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeUserIsNotRegisteredErrContainer()).
				SetComponentsV2(true).
				SetEphemeral(true).
				SetReply(true).
				Send()
		}

		blocked, reason := repository.GetDatabase().Users.IsUserBlocked(ctx, i.User.ID)
		if command.Flags&CommandFlagsIsBlocked != 0 && blocked {
			user, _ := s.User(i.User.ID)
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeUserIsBlockedContainer(user.GlobalName, reason)).
				SetComponentsV2(true).
				SetReply(true).
				Send()
		}

		return command.Run(i)

	}
	return nil
}

func (d *Discommand) ChatInputAutocomplete(name string, s *discordgo.Session, inter *discordgo.InteractionCreate) error {
	i := &builders.InteractionCreate{
		InteractionCreate: inter,
		Session:           s,
	}

	i.InteractionCreate.User = builders.GetInteractionUser(inter)

	if command, ok := d.Commands[name]; ok {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		i.Ctx = ctx

		defer cancel()

		return command.Autocomplete(i)
	}

	return nil
}

func (d *Discommand) ComponentRun(s *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var err error

	i := &builders.InteractionCreate{
		InteractionCreate: inter,
		Session:           s,
	}

	i.InteractionCreate.User = builders.GetInteractionUser(inter)

	for _, c := range d.Components {
		var ctx context.Context
		var cancel context.CancelFunc
		if c.DeferredReply || c.DeferredUpdate {
			ctx, cancel = context.WithTimeout(context.Background(), 15*time.Minute)
		} else {
			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		}
		i.Ctx = ctx

		defer cancel()

		if !c.Parse(i) {
			continue
		}

		if c.DeferredReply {
			if err := i.DeferReply(c.DeferReplyOptions); err != nil {
				return err
			}
		} else if c.DeferredUpdate {
			if err := i.DeferUpdate(); err != nil {
				return err
			}
		}

		err = c.Run(i)
		break
	}
	return err
}

func (d *Discommand) ModalRun(s *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var err error

	i := &builders.InteractionCreate{
		InteractionCreate: inter,
		Session:           s,
	}

	for _, m := range d.Modals {
		var ctx context.Context
		var cancel context.CancelFunc
		if m.Deferred {
			ctx, cancel = context.WithTimeout(context.Background(), 15*time.Minute)
		} else {
			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		}
		i.Ctx = ctx

		defer cancel()

		if !m.Parse(i) {
			continue
		}

		if m.Deferred {
			if err := i.DeferReply(m.DeferOptions); err != nil {
				return err
			}
		}

		err = m.Run(i)
		break
	}
	return err
}
