package commands

import (
	"sync"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/repository"
	"github.com/bwmarrin/discordgo"
)

type run func(inter *builders.InteractionCreate) error
type parse func(inter *builders.InteractionCreate) bool

type CommandFlags uint8

type Command struct {
	*discordgo.ApplicationCommand
	Flags        CommandFlags
	Run          run
	Autocomplete run
}

type Discommand struct {
	Commands   map[string]*Command
	Components []*Component
	Modals     []*Modal
}

type Component struct {
	Parse parse
	Run   run
}

type Modal struct {
	Parse parse
	Run   run
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
		if command.Flags&CommandFlagsIsDeveloperOnlyCommand != 0 && i.User.ID != configs.GetConfig().Bot.OwnerID {
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeDeclineContainer("이 명령어는 개발자 전용 명령어에요.")).
				SetComponentsV2(true).
				SetEphemeral(true).
				Send()
		}

		if command.Flags&CommandFlagsIsRegistered != 0 && !repository.GetDatabase().Users.IsUser(i.User.ID) {
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeUserIsNotRegisteredErrContainer()).
				SetComponentsV2(true).
				SetEphemeral(true).
				SetReply(true).
				Send()
		}

		blocked, reason := repository.GetDatabase().Users.IsUserBlocked(i.User.ID)
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
		if !c.Parse(i) {
			continue
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
		if !m.Parse(i) {
			continue
		}

		err = m.Run(i)
		break
	}
	return err
}
