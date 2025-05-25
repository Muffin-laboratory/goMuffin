package commands

import (
	"context"
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var DeleteLearnedDataCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "삭제",
		Description: "당신이 가르쳐준 단ㅇ어를 삭제해요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "단어",
				Description: "삭제할 단어를 입ㄹ력해주세요.",
				Required:    true,
			},
		},
	},
	Aliases: []string{"잊어", "지워"},
	DetailedDescription: &DetailedDescription{
		Usage:    fmt.Sprintf("%s삭제 (삭제할 단어)", configs.Config.Bot.Prefix),
		Examples: []string{fmt.Sprintf("%s삭제 머핀", configs.Config.Bot.Prefix)},
	},
	Category: Chatting,
	MessageRun: func(ctx *MsgContext) {
		command := strings.Join(*ctx.Args, " ")
		if command == "" {
			utils.NewMessageSender(ctx.Msg).
				AddComponents(utils.GetErrorContainer(
					discordgo.TextDisplay{
						Content: "올바르지 않ㅇ은 용법이에요.",
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **사용법**\n> %s", ctx.Command.DetailedDescription.Usage),
					},
					discordgo.TextDisplay{
						Content: fmt.Sprintf("- **예시**\n%s", strings.Join(utils.AddPrefix("> ", ctx.Command.DetailedDescription.Examples), "\n")),
					},
				)).
				SetComponentsV2(true).
				SetReply(true).
				Send()
		}
		deleteLearnedDataRun(ctx.Msg, strings.Join(*ctx.Args, " "), ctx.Msg.Author.ID)
	},
	ChatInputRun: func(ctx *ChatInputContext) {
		ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		})

		var command string

		if opt, ok := ctx.Inter.Options["단어"]; ok {
			command = opt.StringValue()
		}

		deleteLearnedDataRun(ctx.Inter, command, ctx.Inter.Member.User.ID)
	},
}

func deleteLearnedDataRun(m any, command, userId string) {
	var data []databases.Learn
	var sections []discordgo.Section
	var containers []*discordgo.Container

	cur, err := databases.Database.Learns.Find(context.TODO(), bson.M{"user_id": userId, "command": command})
	if err != nil {
		utils.NewMessageSender(m).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "데이터를 가져오는데 실패했어요."})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
		return
	}

	cur.All(context.TODO(), &data)

	if len(data) < 1 {
		utils.NewMessageSender(m).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "해당 하는 지식ㅇ을 찾을 수 없어요."})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
		return
	}

	for i, data := range data {
		sections = append(sections, discordgo.Section{
			Accessory: discordgo.Button{
				Label:    "삭제",
				Style:    discordgo.DangerButton,
				CustomID: utils.MakeDeleteLearnedData(data.Id.Hex(), i+1, userId),
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

	utils.PaginationEmbedBuilder(m).
		AddContainers(containers...).
		Start()
}
