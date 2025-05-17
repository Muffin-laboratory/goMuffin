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
				AddEmbed(&discordgo.MessageEmbed{
					Title:       "❌ 오류",
					Description: "올바르지 않ㅇ은 용법이에요.",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "사용법",
							Value: utils.InlineCode(ctx.Command.DetailedDescription.Usage),
						},
						{
							Name:  "예시",
							Value: utils.CodeBlock("md", strings.Join(addPrefix(ctx.Command.DetailedDescription.Examples), "\n")),
						},
					},
					Color: utils.EmbedFail,
				}).
				SetReply(true).
				Send()
		}
		deleteLearnedDataRun(ctx.Msg, strings.Join(*ctx.Args, " "), ctx.Msg.Author.ID)
	},
	ChatInputRun: func(ctx *ChatInputContext) {
		ctx.Inter.DeferReply(true)

		var command string

		if opt, ok := ctx.Inter.Options["단어"]; ok {
			command = opt.StringValue()
		}

		deleteLearnedDataRun(ctx.Inter, command, ctx.Inter.Member.User.ID)
	},
}

func deleteLearnedDataRun(m any, command, userId string) {
	var description string
	var data []databases.Learn
	var options []discordgo.SelectMenuOption

	cur, err := databases.Database.Learns.Find(context.TODO(), bson.M{"user_id": userId, "command": command})
	if err != nil {
		utils.NewMessageSender(m).
			AddEmbed(&discordgo.MessageEmbed{
				Title:       "❌ 오류",
				Description: "데이터를 가져오는데 실패했어요.",
				Color:       utils.EmbedFail,
			}).
			SetReply(true).
			Send()
		return
	}

	cur.All(context.TODO(), &data)

	if len(data) < 1 {
		utils.NewMessageSender(m).
			AddEmbed(&discordgo.MessageEmbed{
				Title:       "❌ 오류",
				Description: "해당 하는 지식ㅇ을 찾을 수 없어요.",
				Color:       utils.EmbedFail,
			}).
			SetReply(true).
			Send()
		return
	}

	for i := range len(data) {
		data := data[i]

		options = append(options, discordgo.SelectMenuOption{
			Label:       fmt.Sprintf("%d번 지식", i+1),
			Description: data.Result,
			Value:       utils.MakeDeleteLearnedData(data.Id.Hex(), i+1),
		})
		description += fmt.Sprintf("%d. %s\n", i+1, data.Result)
	}

	utils.NewMessageSender(m).
		AddEmbed(&discordgo.MessageEmbed{
			Title:       fmt.Sprintf("%s 삭제", command),
			Description: utils.CodeBlock("md", fmt.Sprintf("# %s에 대한 대답 중 하나를 선ㅌ택하여 삭제해주세요.\n%s", command, description)),
			Color:       utils.EmbedDefault,
		}).
		AddComponent(discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					MenuType:    discordgo.StringSelectMenu,
					CustomID:    utils.MakeDeleteLearnedDataUserId(userId),
					Options:     options,
					Placeholder: "ㅈ지울 응답을 선택해주세요.",
				},
			},
		}).
		AddComponent(discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: utils.MakeDeleteLearnedDataCancel(userId),
					Label:    "취소하기",
					Style:    discordgo.DangerButton,
					Disabled: false,
				},
			},
		}).
		Send()
}
