package commands

import (
	"sync"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/databases"
	"github.com/bwmarrin/discordgo"
)

type modalRun func(ctx *ModalContext) error
type chatInputRun func(ctx *ChatInputContext) error
type componentRun func(ctx *ComponentContext) error

type modalParse func(ctx *ModalContext) bool
type componentParse func(ctx *ComponentContext) bool

type Category string
type CommandFlags uint8

type Command struct {
	*discordgo.ApplicationCommand
	Category     Category
	Flags        CommandFlags
	Run          chatInputRun
	Autocomplete chatInputRun
}

type Discommand struct {
	Commands   map[string]*Command
	Components []*Component
	Modals     []*Modal
}

type ChatInputContext struct {
	Inter   *builders.InteractionCreate
	Command *Command
}

type ComponentContext struct {
	Inter     *builders.InteractionCreate
	Component *Component
}

type ModalContext struct {
	Inter *builders.InteractionCreate
	Modal *Modal
}

type Component struct {
	Parse componentParse
	Run   componentRun
}

type Modal struct {
	Parse modalParse
	Run   modalRun
}

const (
	Chatting Category = "채팅"
	General  Category = "일반"
)

const (
	CommandFlagsIsRegistered CommandFlags = 1 << iota
	CommandFlagsIsBlocked
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
		if command.Flags&CommandFlagsIsRegistered != 0 && !databases.GetDatabase().Users.IsUser(i.User.ID) {
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeUserIsNotRegisteredErrContainer()).
				SetComponentsV2(true).
				SetEphemeral(true).
				SetReply(true).
				Send()
		}

		blocked, reason := databases.GetDatabase().Users.IsUserBlocked(i.User.ID)
		if command.Flags&CommandFlagsIsBlocked != 0 && blocked {
			user, _ := s.User(i.User.ID)
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeUserIsBlockedContainer(user.GlobalName, reason)).
				SetComponentsV2(true).
				SetReply(true).
				Send()
		}

		return command.Run(&ChatInputContext{i, command})
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
		return command.Autocomplete(&ChatInputContext{i, command})
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
	data := &ComponentContext{
		Inter: i,
	}

	for _, c := range d.Components {
		data.Component = c

		if !c.Parse(data) {
			continue
		}

		err = c.Run(data)
		break
	}
	return err
}

func (d *Discommand) ModalRun(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var err error

	data := &ModalContext{
		Inter: &builders.InteractionCreate{
			InteractionCreate: i,
			Session:           s,
		},
	}

	for _, m := range d.Modals {
		data.Modal = m

		if !m.Parse(data) {
			continue
		}

		err = m.Run(data)
		break
	}
	return err
}
