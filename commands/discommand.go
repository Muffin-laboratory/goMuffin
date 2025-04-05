package commands

import (
	"sync"

	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

type messageRun func(ctx *MsgContext)
type chatInputRun func(ctx *ChatInputContext)
type componentRun func(ctx *ComponentContext)
type componentParse func(ctx *ComponentContext) bool

type Category string

type DetailedDescription struct {
	Usage    string
	Examples []string
}

type Command struct {
	*discordgo.ApplicationCommand
	Aliases             []string
	DetailedDescription *DetailedDescription
	Category            Category
	MessageRun          messageRun
	ChatInputRun        chatInputRun
}

type DiscommandStruct struct {
	Commands   map[string]*Command
	Components []*Component
	Aliases    map[string]string
}

type MsgContext struct {
	Session *discordgo.Session
	Msg     *discordgo.MessageCreate
	Args    []string
	Command *Command
}

type ChatInputContext struct {
	Session *discordgo.Session
	Inter   *utils.InteractionCreate
	Command *Command
}

type ComponentContext struct {
	Session   *discordgo.Session
	Inter     *utils.InteractionCreate
	Component *Component
}

type Component struct {
	Parse componentParse
	Run   componentRun
}

const (
	Chattings Category = "채팅"
	Generals  Category = "일반"
)

var commandMutex *sync.Mutex = &sync.Mutex{}
var componentMutex *sync.Mutex = &sync.Mutex{}

func new() *DiscommandStruct {
	discommand := DiscommandStruct{
		Commands:   map[string]*Command{},
		Aliases:    map[string]string{},
		Components: []*Component{},
	}
	return &discommand
}

func (d *DiscommandStruct) LoadCommand(c *Command) {
	commandMutex.Lock()
	d.Commands[c.Name] = c
	d.Aliases[c.Name] = c.Name

	for _, alias := range c.Aliases {
		d.Aliases[alias] = c.Name
	}
	commandMutex.Unlock()
}

func (d *DiscommandStruct) LoadComponent(c *Component) {
	componentMutex.Lock()
	d.Components = append(d.Components, c)
	componentMutex.Unlock()
}

func (d *DiscommandStruct) MessageRun(name string, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if command, ok := d.Commands[name]; ok {
		command.MessageRun(&MsgContext{s, m, args, command})
	}
}

func (d *DiscommandStruct) ChatInputRun(name string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if command, ok := d.Commands[name]; ok {
		command.ChatInputRun(&ChatInputContext{s, &utils.InteractionCreate{
			InteractionCreate: i,
			Session:           s,
			Options:           utils.GetInteractionOptions(i),
		}, command})
	}
}

func (d *DiscommandStruct) ComponentRun(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, c := range d.Components {
		if (!c.Parse(&ComponentContext{s, &utils.InteractionCreate{
			InteractionCreate: i,
			Session:           s,
		}, c})) {
			continue
		}

		c.Run(&ComponentContext{s, &utils.InteractionCreate{
			InteractionCreate: i,
			Session:           s,
		}, c})
		break
	}
}

var Discommand *DiscommandStruct = new()
