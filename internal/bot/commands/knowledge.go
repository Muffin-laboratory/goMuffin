package commands

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	subcommands "github.com/Muffin-laboratory/goMuffin/internal/bot/commands/subcommands/knowledge"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	knowledgeLearn  = "추가"
	knowledgeList   = "목록"
	knowledgeDelete = "삭제"
)

var KnowledgeCommand = &loader.Command{
	Deferred: true,
	DeferOptions: &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
	},
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
	Flags: loader.CommandFlagsIsRegistered | loader.CommandFlagsIsBlocked,
	Run: func(inter *builders.InteractionCreate) error {
		switch opt := inter.ApplicationCommandData().Options[0]; opt.Name {
		case knowledgeLearn:
			igCommands := []string{}

			for _, command := range loader.GetDiscommand().Commands {
				igCommands = append(igCommands, command.Name)
			}

			return subcommands.Learn(inter, builders.MakeCommandInteractionOptionsMap(opt.Options), igCommands)
		case knowledgeList:
			return subcommands.List(inter, builders.MakeCommandInteractionOptionsMap(opt.Options))
		case knowledgeDelete:
			return subcommands.Delete(inter, builders.MakeCommandInteractionOptionsMap(opt.Options))
		default:
			return nil
		}
	},
	Autocomplete: func(inter *builders.InteractionCreate) error {
		var choices []*discordgo.ApplicationCommandOptionChoice
		var focusedValue string

		for _, opt := range inter.ApplicationCommandData().Options[0].Options {
			if opt.Focused {
				focusedValue = opt.StringValue()
				break
			}
		}

		data, err := repository.GetDatabase().Knowledge.GetByFilter(inter.Ctx, bson.M{
			"user_id": inter.User.ID,
			"command": bson.M{
				"$regex": focusedValue,
			},
		})
		if err != nil {
			return err
		}

		for _, data := range data {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  data.Command,
				Value: data.Command,
			})
		}

		return inter.Autocomplete(choices)
	},
}

func init() {
	loader.GetDiscommand().LoadCommand(KnowledgeCommand)
}
