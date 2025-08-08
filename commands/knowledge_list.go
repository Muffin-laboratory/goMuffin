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

var KnowledgeListCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "리스트",
		Description: "당신이 가ㄹ르쳐준 지식을 나열해요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "단어",
				Description: "해당 단어에 대한 결과를 찾아요.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "개수",
				Description: "한 페이지당 보여줄 지식 데이터 양을 정해요.",
				MinValue:    &LIST_MIN_VALUE,
				MaxValue:    LIST_MAX_VALUE,
				Required:    false,
			},
		},
	},
	Aliases: []string{"list", "목록", "지식목록"},
	DetailedDescription: &DetailedDescription{
		Usage: configs.AddPrefix("%s리스트 [단어]"),
		Examples: []string{
			configs.AddPrefix("%s리스트"),
			configs.AddPrefix("%s리스트 안녕"),
			configs.AddPrefix("%s리스트 개수:10"),
		},
	},
	Category:                   Chatting,
	RegisterApplicationCommand: true,
	RegisterMessageCommand:     true,
	Flags:                      CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	MessageRun: func(ctx *MsgContext) error {
		var length int

		filter := bson.D{{Key: "user_id", Value: ctx.Msg.Author.ID}}
		query := strings.Join(*ctx.Args, " ")

		command := utils.RegexpLearnQueryLength.ReplaceAllString(query, "")
		command = strings.Join(strings.Fields(command), " ")
		if command != "" {
			filter = append(filter, bson.E{
				Key:   "command",
				Value: command,
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
				return nil
			}

			if float64(length) > LIST_MAX_VALUE {
				utils.NewMessageSender(ctx.Msg).
					AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: fmt.Sprintf("개수의 값은 %d보다 작아야해요.", int(LIST_MAX_VALUE))})).
					SetComponentsV2(true).
					SetReply(true).
					Send()
				return nil
			}
		}
		return learnedDataListRun(ctx.Msg, ctx.Msg.Author.GlobalName, ctx.Msg.Author.AvatarURL("512"), filter, length)
	},
	ChatInputRun: func(ctx *ChatInputContext) error {
		err := ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			return err
		}

		var length int

		filter := bson.D{{Key: "user_id", Value: ctx.Inter.Member.User.ID}}

		if opt, ok := ctx.Inter.Options["단어"]; ok {
			filter = append(filter, bson.E{
				Key:   "command",
				Value: opt.StringValue(),
			})
		}

		if opt, ok := ctx.Inter.Options["개수"]; ok {
			length = int(opt.IntValue())
		}
		return learnedDataListRun(ctx.Inter, ctx.Inter.Member.User.GlobalName, ctx.Inter.Member.User.AvatarURL("512"), filter, length)
	},
}

func getDescriptions(items []string, length int) (descriptions []string) {
	var builder strings.Builder

	if length == 0 {
		length = 25
	}

	for i, item := range items {
		builder.WriteString(fmt.Sprintf("%s\n", item))

		if (i+1)%length == 0 {
			descriptions = append(descriptions, builder.String())
			builder.Reset()
		}
		i += 1
	}

	if builder.Len() > 0 {
		descriptions = append(descriptions, builder.String())
	}
	return
}

func getContainers(accessory *discordgo.Thumbnail, defaultDesc string, items []string, length int) []*discordgo.Container {
	var containers []*discordgo.Container

	descriptions := getDescriptions(items, length)

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

func learnedDataListRun(m any, globalName, avatarURL string, filter bson.D, length int) error {
	var data []databases.Learn

	itemsMap := map[string]string{}
	items := []string{}

	cur, err := databases.GetDatabase().Learns.Find(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.NewMessageSender(m).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "당신은 지식ㅇ을 가르쳐준 적이 없어요!"})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return nil
		}

		utils.NewMessageSender(m).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "데이터를 가져오는데 실패했어요."})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
		return err
	}

	defer cur.Close(context.TODO())

	cur.All(context.TODO(), &data)

	if len(filter) > 1 {
		command := filter[1].Value.(string)

		for _, data := range data {
			items = append(items, fmt.Sprintf("> %s", data.Result))
		}

		containers := getContainers(&discordgo.Thumbnail{
			Media: discordgo.UnfurledMediaItem{
				URL: avatarURL,
			},
		}, fmt.Sprintf("### %s님이 알려주신 지식\n- **%s**\n", globalName, command)+"%s", items, length)

		return utils.PaginationContainerBuilder(m).
			AddContainers(containers...).
			Start()
	}

	for _, data := range data {
		if _, ok := itemsMap[data.Command]; ok {
			continue
		}

		itemsMap[data.Command] = fmt.Sprintf("- `%s`", data.Command)
	}

	for _, v := range itemsMap {
		items = append(items, v)
	}

	containers := getContainers(&discordgo.Thumbnail{
		Media: discordgo.UnfurledMediaItem{
			URL: avatarURL,
		},
	}, fmt.Sprintf("### %s님이 알려주신 지식\n총 %d개에요.\n", globalName, len(items))+"%s", items, length)

	return utils.PaginationContainerBuilder(m).
		AddContainers(containers...).
		Start()
}
