package commands

import (
	"sync"

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

func (d *DiscommandStruct) MessageRun(name string, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if command, ok := d.Commands[name]; ok {
		command.MessageRun(&MsgContext{&utils.MessageCreate{
			MessageCreate: m,
			Session:       s,
		}, &args, command})
	}
}

func (d *DiscommandStruct) ChatInputRun(name string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if command, ok := d.Commands[name]; ok {
		command.ChatInputRun(&ChatInputContext{&utils.InteractionCreate{
			InteractionCreate: i,
			Session:           s,
			Options:           utils.GetInteractionOptions(i),
		}, command})
	}
}

func (d *DiscommandStruct) ComponentRun(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := &ComponentContext{
		Inter: &utils.InteractionCreate{
			InteractionCreate: i,
			Session:           s,
		},
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
