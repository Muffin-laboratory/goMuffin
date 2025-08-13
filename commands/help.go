package commands

import (
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var HelpCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "도움말",
		Description: "기본적인 사용ㅂ법이에요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "명령어",
				Description: "도움말을 볼 명령어",
				Choices:     []*discordgo.ApplicationCommandOptionChoice{},
			},
		},
	},
	Aliases: []string{"도움", "명령어", "help"},
	DetailedDescription: &DetailedDescription{
		Usage:    "/도움말 [명령어:문자(머핀봇의 명령어)]",
		Examples: []string{"/도움말", "/도움말 명령어:배워"},
	},
	Category:                   General,
	RegisterApplicationCommand: true,
	RegisterMessageCommand:     true,
	Flags:                      CommandFlagsIsBlocked,
	MessageRun: func(ctx *MsgContext) error {
		return helpRun(ctx.Msg.Session, ctx.Msg, strings.Join(*ctx.Args, " "))
	},
	ChatInputRun: func(ctx *ChatInputContext) error {
		var command string

		if opt, ok := ctx.Inter.Options["명령어"]; ok {
			command = opt.StringValue()
		}

		return helpRun(ctx.Inter.Session, ctx.Inter, command)
	},
}

func getCommandsByCategory(d *DiscommandStruct, category Category) []string {
	commands := []string{}
	for _, command := range d.Commands {
		if command.Category == category {
			commands = append(commands, fmt.Sprintf("> **%s**: %s", command.Name, command.Description))
		}
	}
	return commands
}

func helpRun(s *discordgo.Session, m any, commandName string) error {
	section := &discordgo.Section{
		Accessory: discordgo.Thumbnail{
			Media: discordgo.UnfurledMediaItem{
				URL: s.State.User.AvatarURL("512"),
			},
		},
	}

	commandName = Discommand.Aliases[commandName]

	if commandName == "" || Discommand.Commands[commandName] == nil {
		section.Components = append(section.Components,
			discordgo.TextDisplay{
				Content: fmt.Sprintf("### %s의 도움말", s.State.User.Username),
			},
			discordgo.TextDisplay{
				Content: fmt.Sprintf("- **일반**\n%s", strings.Join(getCommandsByCategory(Discommand, General), "\n")),
			},
			discordgo.TextDisplay{
				Content: fmt.Sprintf("- **채팅**\n%s", strings.Join(getCommandsByCategory(Discommand, Chatting), "\n")),
			},
		)
		return utils.NewMessageSender(m).
			AddComponents(&discordgo.Container{
				Components: []discordgo.MessageComponent{
					section,
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label: "개인정보처리방침",
								URL:   configs.Config.Service.PrivacyPolicyURL,
								Style: discordgo.LinkButton,
								Emoji: &discordgo.ComponentEmoji{
									Name: "🔗",
								},
							},
							discordgo.Button{
								Label: "서비스 이용약관",
								URL:   configs.Config.Service.TermOfServiceURL,
								Style: discordgo.LinkButton,
								Emoji: &discordgo.ComponentEmoji{
									Name: "🔗",
								},
							},
						},
					},
				},
			}).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	}

	var aliases, examples discordgo.TextDisplay

	command := Discommand.Commands[commandName]

	section.Components = append(section.Components,
		discordgo.TextDisplay{
			Content: fmt.Sprintf("### %s의 %s 명령어의 도움말", s.State.User.Username, command.Name),
		},
		discordgo.TextDisplay{
			Content: fmt.Sprintf("- **설명**\n> %s", command.Description),
		},
		discordgo.TextDisplay{
			Content: fmt.Sprintf("- **사용법**\n> %s", command.DetailedDescription.Usage),
		},
	)

	if command.Aliases != nil {
		aliases = discordgo.TextDisplay{
			Content: fmt.Sprintf("- **별칭**\n%s", strings.Join(utils.AddPrefix("> ", command.Aliases), "\n")),
		}
	} else {
		aliases = discordgo.TextDisplay{
			Content: "- **별칭**\n> 없음",
		}
	}

	if command.DetailedDescription.Examples != nil {
		examples = discordgo.TextDisplay{
			Content: fmt.Sprintf("- **예시**\n%s", strings.Join(utils.AddPrefix("> ", command.DetailedDescription.Examples), "\n")),
		}
	} else {
		aliases = discordgo.TextDisplay{
			Content: "- **예시**\n> 없음",
		}
	}

	if command.Name == LearnCommand.Name {
		learnArgs := discordgo.TextDisplay{
			Content: fmt.Sprintf("- **대답에 쓸 수 있는 인자**\n%s", learnArguments),
		}
		return utils.NewMessageSender(m).
			AddComponents(discordgo.Container{
				Components: []discordgo.MessageComponent{section, aliases, examples, learnArgs},
			}).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	}

	return utils.NewMessageSender(m).
		AddComponents(discordgo.Container{
			Components: []discordgo.MessageComponent{section, aliases, examples},
		}).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
