package commands

import (
	subcommands "git.wh64.net/muffin/goMuffin/commands/subcommands/knowledge"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

const (
	knowledgeLearn  = "추가"
	knowledgeList   = "리스트"
	knowledgeDelete = "삭제"
)

var KnowledgeCommand = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "지식",
		Description: "이 봇이 사용자와 대화할 때 알면 좋은 지식을 관리하는 명령어에요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        knowledgeLearn,
				Description: "단어를 가르치는 명령어에요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "단어",
						Description: "등록할 단어",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "대답",
						Description: "해당 단어의 대답",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        knowledgeList,
				Description: "당신이 가르쳐준 지식을 나열해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "단어",
						Description: "해당 단어가 포함된 결과",
						Required:    false,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        chatCommandDelete,
				Description: "당신이 가르쳐준 단어를 삭제해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "단어",
						Description: "삭제할 단어",
						Required:    true,
					},
				},
			},
		},
	},
	Category: Chatting,
	Flags:    CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	Run: func(ctx *ChatInputContext) error {
		switch opt := ctx.Inter.ApplicationCommandData().Options[0]; opt.Name {
		case knowledgeLearn:
			if err := ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			igCommands := []string{}

			for _, command := range instance.Commands {
				igCommands = append(igCommands, command.Name)
			}

			return subcommands.Learn(ctx.Inter, utils.MakeCommandInteractionOptionsMap(opt.Options), igCommands)
		case knowledgeList:
			return subcommands.List(ctx.Inter)
		case knowledgeDelete:
			return subcommands.Delete(ctx.Inter)
		default:
			return nil
		}
	},
}

func init() {
	GetDiscommand().LoadCommand(KnowledgeCommand)
}
