package commands

import (
	subcommands "git.wh64.net/muffin/goMuffin/commands/subcommands/knowledge"
	"github.com/bwmarrin/discordgo"
)

const (
	knowledgeLearn  = "배워"
	knowledgeList   = "리스트"
	knowledgeDelete = "삭제"
)

var KnowledgeCommand = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name: "지식",
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
			return subcommands.Learn(ctx.Inter)
		case knowledgeList:
			return subcommands.Learn(ctx.Inter)
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
