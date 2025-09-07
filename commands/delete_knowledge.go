package commands

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var DeleteKnowledgeCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "삭제",
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
	DetailedDescription: DetailedDescription{
		Usage:    "/삭제 (단어:문자)",
		Examples: []string{"/삭제 단어:뷁"},
	},
	Category: Chatting,
	Flags:    CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	Run: func(ctx *ChatInputContext) error {
		err := ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			return err
		}

		var command string

		if opt, ok := ctx.Inter.Options["단어"]; ok {
			command = opt.StringValue()
		}

		return deleteLearnedDataRun(ctx.Inter, command, ctx.Inter.Member.User.ID)
	},
}

func deleteLearnedDataRun(m any, command, userID string) error {
	var data []databases.Learn
	var sections []discordgo.Section
	var containers []*discordgo.Container

	cur, err := databases.GetDatabase().Learns.Find(context.TODO(), databases.Learn{UserID: userID, Command: command})
	if err != nil {
		return err
	}

	cur.All(context.TODO(), &data)

	if len(data) < 1 {
		utils.NewMessageSender(m).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "해당 하는 지식ㅇ을 찾을 수 없어요."})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
		return nil
	}

	for i, data := range data {
		sections = append(sections, discordgo.Section{
			Accessory: discordgo.Button{
				Label:    "삭제",
				Style:    discordgo.DangerButton,
				CustomID: utils.MakeDeleteLearnedData(data.ID.Hex(), i+1, userID),
			},
			Components: []discordgo.MessageComponent{
				discordgo.TextDisplay{
					Content: fmt.Sprintf("%d. %s\n", i+1, data.Result),
				},
			},
		})
	}

	textDisplay := discordgo.TextDisplay{Content: fmt.Sprintf("### %s 삭제", command)}
	container := &discordgo.Container{Components: []discordgo.MessageComponent{textDisplay}}
	for i, section := range sections {
		container.Components = append(container.Components, section, discordgo.Separator{})

		if (i+1)%10 == 0 {
			containers = append(containers, container)
			container = &discordgo.Container{Components: []discordgo.MessageComponent{textDisplay}}
			continue
		}
	}

	if len(container.Components) > 1 {
		containers = append(containers, container)
	}

	return utils.PaginationContainerBuilder(m).
		AddContainers(containers...).
		Start()
}

func init() {
	GetDiscommand().LoadCommand(DeleteKnowledgeCommand)
}
