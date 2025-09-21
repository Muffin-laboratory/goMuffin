package commands

import (
	"context"

	subcommands "git.wh64.net/muffin/goMuffin/commands/subcommands/knowledge"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	knowledgeLearn  = "추가"
	knowledgeList   = "목록"
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
						MaxLength:   100,
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
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "단어",
						Description:  "해당 단어가 포함된 결과",
						Required:     false,
						Autocomplete: true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        chatCommandDelete,
				Description: "당신이 가르쳐준 단어를 삭제해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "단어",
						Description:  "삭제할 단어",
						Required:     true,
						Autocomplete: true,
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
			if err := ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			return subcommands.List(ctx.Inter, utils.MakeCommandInteractionOptionsMap(opt.Options))
		case knowledgeDelete:
			if err := ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			}); err != nil {
				return err
			}

			return subcommands.Delete(ctx.Inter, utils.MakeCommandInteractionOptionsMap(opt.Options))
		default:
			return nil
		}
	},
	Autocomplete: func(ctx *ChatInputContext) error {
		var choices []*discordgo.ApplicationCommandOptionChoice
		var focusedValue string
		var data []*databases.Knowledge

		for _, opt := range ctx.Inter.ApplicationCommandData().Options[0].Options {
			if opt.Focused {
				focusedValue = opt.StringValue()
				break
			}
		}

		cur, err := databases.GetDatabase().Knowledge.Find(context.TODO(), bson.M{
			"user_id": ctx.Inter.User.ID,
			"command": bson.M{
				"$regex": focusedValue,
			},
		})
		if err != nil {
			return err
		}

		defer cur.Close(context.TODO())

		if err = cur.All(context.TODO(), &data); err != nil {
			return err
		}

		for _, data := range data {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  data.Command,
				Value: data.Command,
			})
		}

		return ctx.Inter.Autocomplete(choices)
	},
}

func init() {
	GetDiscommand().LoadCommand(KnowledgeCommand)
}
