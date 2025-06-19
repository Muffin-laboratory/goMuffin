package commands

import (
	"sync"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

type modalRun func(ctx *ModalContext)
type messageRun func(ctx *MsgContext)
type chatInputRun func(ctx *ChatInputContext)
type componentRun func(ctx *ComponentContext)

type modalParse func(ctx *ModalContext) bool
type componentParse func(ctx *ComponentContext) bool

type Category string
type CommandFlags uint8

type DetailedDescription struct {
	Usage    string
	Examples []string
}

type Command struct {
	*discordgo.ApplicationCommand
	Aliases                    []string
	DetailedDescription        *DetailedDescription
	Category                   Category
	RegisterApplicationCommand bool
	RegisterMessageCommand     bool
	Flags                      CommandFlags
	MessageRun                 messageRun
	ChatInputRun               chatInputRun
}

type DiscommandStruct struct {
	Commands   map[string]*Command
	Components []*Component
	Aliases    map[string]string
	Modals     []*Modal
}

type MsgContext struct {
	Msg     *utils.MessageCreate
	Args    *[]string
	Command *Command
}

type ChatInputContext struct {
	Inter   *utils.InteractionCreate
	Command *Command
}

type ComponentContext struct {
	Inter     *utils.InteractionCreate
	Component *Component
}

type ModalContext struct {
	Inter *utils.InteractionCreate
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
	Chatting      Category = "채팅"
	General       Category = "일반"
	DeveloperOnly Category = "개발자 전용"
)

const (
	CommandFlagsIsDeveloper CommandFlags = 1 << iota
	CommandFlagsIsRegistered
)

var (
	commandMutex   sync.Mutex
	componentMutex sync.Mutex
	modalMutex     sync.Mutex
)

var Discommand *DiscommandStruct

func init() {
	Discommand = &DiscommandStruct{
		Commands:   map[string]*Command{},
		Aliases:    map[string]string{},
		Components: []*Component{},
		Modals:     []*Modal{},
	}
}

func (d *DiscommandStruct) LoadCommand(c *Command) {
	defer commandMutex.Unlock()
	commandMutex.Lock()
	d.Commands[c.Name] = c
	d.Aliases[c.Name] = c.Name

	for _, alias := range c.Aliases {
		d.Aliases[alias] = c.Name
	}
}

func (d *DiscommandStruct) LoadComponent(c *Component) {
	defer componentMutex.Unlock()
	componentMutex.Lock()
	d.Components = append(d.Components, c)
}

func (d *DiscommandStruct) LoadModal(m *Modal) {
	defer modalMutex.Unlock()
	modalMutex.Lock()
	d.Modals = append(d.Modals, m)
}

func (d *DiscommandStruct) MessageRun(name string, s *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	m := &utils.MessageCreate{
		MessageCreate: msg,
		Session:       s,
	}

	if command, ok := d.Commands[name]; ok && command.RegisterMessageCommand {
		if command.Flags&CommandFlagsIsDeveloper != 0 && m.Author.ID != configs.Config.Bot.OwnerId {
			utils.NewMessageSender(m).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "해당 명령어는 개발자만 사용 가능해요."})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return
		}

		if command.Flags&CommandFlagsIsRegistered != 0 && !databases.Database.IsUser(m.Author.ID) {
			utils.NewMessageSender(m).
				AddComponents(utils.GetUserIsNotRegisteredErrContainer(configs.Config.Bot.Prefix)).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return
		}

		command.MessageRun(&MsgContext{m, &args, command})
	}
}

func (d *DiscommandStruct) ChatInputRun(name string, s *discordgo.Session, inter *discordgo.InteractionCreate) {
	i := &utils.InteractionCreate{
		InteractionCreate: inter,
		Session:           s,
		Options:           utils.GetInteractionOptions(inter),
	}

	i.InteractionCreate.User = utils.GetInteractionUser(inter)

	if command, ok := d.Commands[name]; ok && command.RegisterApplicationCommand {
		if command.Flags&CommandFlagsIsDeveloper != 0 && i.User.ID != configs.Config.Bot.OwnerId {
			utils.NewMessageSender(i).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "해당 명령어는 개발자만 사용 가능해요."})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return
		}

		if command.Flags&CommandFlagsIsRegistered != 0 && !databases.Database.IsUser(i.User.ID) {
			utils.NewMessageSender(i).
				AddComponents(utils.GetUserIsNotRegisteredErrContainer(configs.Config.Bot.Prefix)).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return
		}

		command.ChatInputRun(&ChatInputContext{i, command})
	}
}

func (d *DiscommandStruct) ComponentRun(s *discordgo.Session, inter *discordgo.InteractionCreate) {
	i := &utils.InteractionCreate{
		InteractionCreate: inter,
		Session:           s,
	}

	i.InteractionCreate.User = utils.GetInteractionUser(inter)
	data := &ComponentContext{
		Inter: i,
	}

	for _, c := range d.Components {
		data.Component = c

		if !c.Parse(data) {
			continue
		}

		c.Run(data)
		break
	}
}

func (d *DiscommandStruct) ModalRun(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := &ModalContext{
		Inter: &utils.InteractionCreate{
			InteractionCreate: i,
			Session:           s,
		},
	}

	for _, m := range d.Modals {
		data.Modal = m

		if !m.Parse(data) {
			continue
		}

		m.Run(data)
		break
	}
}
