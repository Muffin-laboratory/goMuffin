package commands

import (
	"sync"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

type modalRun func(ctx *ModalContext) error
type messageRun func(ctx *MsgContext) error
type chatInputRun func(ctx *ChatInputContext) error
type componentRun func(ctx *ComponentContext) error

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
	CommandFlagsIsBlocked
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

func (d *DiscommandStruct) MessageRun(name string, s *discordgo.Session, msg *discordgo.MessageCreate, args []string) error {
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
			return nil
		}

		if command.Flags&CommandFlagsIsRegistered != 0 && !databases.Database.IsUser(m.Author.ID) {
			utils.NewMessageSender(m).
				AddComponents(utils.GetUserIsNotRegisteredErrContainer(configs.Config.Bot.Prefix)).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return nil
		}

		blocked, reason := databases.Database.IsUserBlocked(m.Author.ID)
		if command.Flags&CommandFlagsIsBlocked != 0 && blocked {
			user, _ := s.User(m.Author.ID)
			utils.NewMessageSender(m).
				AddComponents(utils.GetUserIsBlockedContainer(user.GlobalName, reason)).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return nil
		}

		utils.NewMessageSender(m).
			AddComponents(discordgo.Container{
				Components: []discordgo.MessageComponent{
					discordgo.TextDisplay{Content: "### ⚠️ 고지"},
					discordgo.TextDisplay{
						Content: "메세지 기반 명령어는 머핀봇 7.0.0 (MadeleineV2)부터 지원이 종료될 예정이에요. " +
							"따라서 앞으로는 빗금 기반 명령어를 사용해주세요.",
					},
				},
			}).
			SetReply(true).
			SetComponentsV2(true).
			Send()

		return command.MessageRun(&MsgContext{m, &args, command})
	}
	return nil
}

func (d *DiscommandStruct) ChatInputRun(name string, s *discordgo.Session, inter *discordgo.InteractionCreate) error {
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
				SetEphemeral(true).
				SetReply(true).
				Send()
			return nil
		}

		if command.Flags&CommandFlagsIsRegistered != 0 && !databases.Database.IsUser(i.User.ID) {
			utils.NewMessageSender(i).
				AddComponents(utils.GetUserIsNotRegisteredErrContainer(configs.Config.Bot.Prefix)).
				SetComponentsV2(true).
				SetEphemeral(true).
				SetReply(true).
				Send()
			return nil
		}

		blocked, reason := databases.Database.IsUserBlocked(i.User.ID)
		if command.Flags&CommandFlagsIsBlocked != 0 && blocked {
			user, _ := s.User(i.User.ID)
			utils.NewMessageSender(i).
				AddComponents(utils.GetUserIsBlockedContainer(user.GlobalName, reason)).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return nil
		}

		return command.ChatInputRun(&ChatInputContext{i, command})
	}
	return nil
}

func (d *DiscommandStruct) ComponentRun(s *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var err error

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

		err = c.Run(data)
		break
	}
	return err
}

func (d *DiscommandStruct) ModalRun(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var err error

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

		err = m.Run(data)
		break
	}
	return err
}
