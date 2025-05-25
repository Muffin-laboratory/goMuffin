package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	LIST_MIN_VALUE float64 = 10.0
	LIST_MAX_VALUE float64 = 100.0
)

var LearnedDataListCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "리스트",
		Description: "당신이 가ㄹ르쳐준 지식을 나열해요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "단어",
				Description: "해당 단어가 포함된 결과만 찾아요.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "대답",
				Description: "해당 대답이 포함된 결과만 찾아요.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "개수",
				Description: "한 페이지당 보여줄 지식의 데이터 양을 정해요.",
				MinValue:    &LIST_MIN_VALUE,
				MaxValue:    LIST_MAX_VALUE,
				Required:    false,
			},
		},
	},
	Aliases: []string{"list", "목록", "지식목록"},
	DetailedDescription: &DetailedDescription{
		Usage: fmt.Sprintf("%s리스트", configs.Config.Bot.Prefix),
		Examples: []string{
			fmt.Sprintf("%s리스트 ㅁㄴㅇㄹ", configs.Config.Bot.Prefix),
			fmt.Sprintf("%s리스트 단어:안녕", configs.Config.Bot.Prefix),
			fmt.Sprintf("%s리스트 대답:머핀", configs.Config.Bot.Prefix),
		},
	},
	Category: Chatting,
	MessageRun: func(ctx *MsgContext) {
		var length int

		filter := bson.D{{Key: "user_id", Value: ctx.Msg.Author.ID}}
		query := strings.Join(*ctx.Args, " ")

		if match := utils.RegexpLearnQueryCommand.FindStringSubmatch(query); match != nil {
			filter = append(filter, bson.E{
				Key: "command",
				Value: bson.M{
					"$regex": match[1],
				},
			})
		}

		if match := utils.RegexpLearnQueryResult.FindStringSubmatch(query); match != nil {
			filter = append(filter, bson.E{
				Key: "result",
				Value: bson.M{
					"$regex": match[1],
				},
			})
		}

		if match := utils.RegexpLearnQueryLength.FindStringSubmatch(query); match != nil {
			length, _ = strconv.Atoi(match[1])

			if float64(length) < LIST_MIN_VALUE {
				utils.NewMessageSender(ctx.Msg).
					AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: fmt.Sprintf("개수의 값은 %d보다 커야해요.", int(LIST_MIN_VALUE))})).
					SetComponentsV2(true).
					SetReply(true).
					Send()
				return
			}

			if float64(length) > LIST_MAX_VALUE {
				utils.NewMessageSender(ctx.Msg).
					AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: fmt.Sprintf("개수의 값은 %d보다 작아야해요.", int(LIST_MAX_VALUE))})).
					SetComponentsV2(true).
					SetReply(true).
					Send()
				return
			}
		}
		learnedDataListRun(ctx.Msg, ctx.Msg.Author.GlobalName, ctx.Msg.Author.AvatarURL("512"), filter, length)
	},
	ChatInputRun: func(ctx *ChatInputContext) {
		ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		})

		var length int

		filter := bson.D{{Key: "user_id", Value: ctx.Inter.Member.User.ID}}

		if opt, ok := ctx.Inter.Options["단어"]; ok {
			filter = append(filter, bson.E{
				Key: "command",
				Value: bson.M{
					"$regex": opt.StringValue(),
				},
			})
		}

		if opt, ok := ctx.Inter.Options["대답"]; ok {
			filter = append(filter, bson.E{
				Key: "result",
				Value: bson.M{
					"$regex": opt.StringValue(),
				},
			})
		}

		if opt, ok := ctx.Inter.Options["개수"]; ok {
			length = int(opt.IntValue())
		}
		learnedDataListRun(ctx.Inter, ctx.Inter.Member.User.GlobalName, ctx.Inter.Member.User.AvatarURL("512"), filter, length)
	},
}

func getDescriptions(data *[]databases.Learn, length int) (descriptions []string) {
	var builder strings.Builder
	MAX_LENGTH := 100

	if length == 0 {
		length = 25
	}

	tempDesc := []string{}

	for _, data := range *data {
		command := data.Command
		result := data.Result

		if runeCommand := []rune(command); len(runeCommand) >= MAX_LENGTH {
			command = string(runeCommand)[:MAX_LENGTH] + "..."
		}

		if runeResult := []rune(result); len(runeResult) >= MAX_LENGTH {
			result = string(runeResult[:MAX_LENGTH]) + "..."
		}

		tempDesc = append(tempDesc, fmt.Sprintf("- %s: %s\n", command, result))
	}

	for i, s := range tempDesc {
		builder.WriteString(s)

		if (i+1)%length == 0 {
			descriptions = append(descriptions, builder.String())
			builder.Reset()
		}
	}

	if builder.Len() > 0 {
		descriptions = append(descriptions, builder.String())
	}
	return
}

func getContainers(accessory *discordgo.Thumbnail, defaultDesc string, data *[]databases.Learn, length int) []*discordgo.Container {
	var containers []*discordgo.Container

	descriptions := getDescriptions(data, length)

	if len(descriptions) <= 0 {
		containers = append(containers, &discordgo.Container{
			Components: []discordgo.MessageComponent{
				discordgo.Section{
					Accessory: accessory,
					Components: []discordgo.MessageComponent{
						discordgo.TextDisplay{
							Content: utils.MakeDesc(defaultDesc, "없음"),
						},
					},
				},
			},
		})
	}

	for _, desc := range descriptions {
		containers = append(containers, &discordgo.Container{
			Components: []discordgo.MessageComponent{
				discordgo.Section{
					Accessory: accessory,
					Components: []discordgo.MessageComponent{
						discordgo.TextDisplay{
							Content: utils.MakeDesc(defaultDesc, desc),
						},
					},
				},
			},
		})
	}
	return containers
}

func learnedDataListRun(m any, globalName, avatarUrl string, filter bson.D, length int) {
	var data []databases.Learn

	cur, err := databases.Database.Learns.Find(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.NewMessageSender(m).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "당신은 지식ㅇ을 가르쳐준 적이 없어요!"})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return
		}

		fmt.Println(err)

		utils.NewMessageSender(m).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "데이터를 가져오는데 실패했어요."})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
		return
	}

	defer cur.Close(context.TODO())

	cur.All(context.TODO(), &data)

	containers := getContainers(&discordgo.Thumbnail{
		Media: discordgo.UnfurledMediaItem{
			URL: avatarUrl,
		},
	}, fmt.Sprintf("### %s님이 알려주신 지식\n%s", globalName, utils.CodeBlock("md", fmt.Sprintf("# 총 %d개에요.\n", len(data))+"%s")), &data, length)

	utils.PaginationEmbedBuilder(m).
		AddContainers(containers...).
		Start()
}
