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
		learnedDataListRun(ctx.Session, ctx.Msg, ctx.Args)
	},
	ChatInputRun: func(ctx *ChatInputContext) {
		learnedDataListRun(ctx.Session, ctx.Inter, nil)
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

func learnedDataListRun(s *discordgo.Session, m any, args *[]string) {
	var globalName, avatarUrl string
	var data []databases.Learn
	var filter bson.D
	var length int

	switch m := m.(type) {
	case *discordgo.MessageCreate:
		filter = bson.D{{Key: "user_id", Value: m.Author.ID}}
		globalName = m.Author.GlobalName
		avatarUrl = m.Author.AvatarURL("512")

		query := strings.Join(*args, " ")

		if match := utils.RegexpLearnQueryCommand.FindStringSubmatch(query); match != nil {
			filter = append(filter, bson.E{
				Key: "command",
				Value: bson.M{
					"$regex": match[1],
				},
			})
		}

		if match := utils.RegexpLearnQueryResult.FindStringSubmatch(query); match != nil {
			fmt.Println(match[1])
			filter = append(filter, bson.E{
				Key: "result",
				Value: bson.M{
					"$regex": match[1],
				},
			})
		}

		if match := utils.RegexpLearnQueryLength.FindStringSubmatch(query); match != nil {
			fmt.Println(1)
			var err error
			length, err = strconv.Atoi(match[1])
			fmt.Printf("err: %v\n", err)

			if err != nil {
				s.ChannelMessageSendEmbedReply(m.ChannelID, &discordgo.MessageEmbed{
					Title:       "❌ 오류",
					Description: "개수의 값은 숫자여야해요.",
					Color:       utils.EmbedFail,
				}, m.Reference())
				return
			}
		}
	case *utils.InteractionCreate:
		m.DeferReply(true)

		filter = bson.D{{Key: "user_id", Value: m.Member.User.ID}}
		globalName = m.Member.User.GlobalName
		avatarUrl = m.Member.User.AvatarURL("512")

		if opt, ok := m.Options["단어"]; ok {
			filter = append(filter, bson.E{
				Key: "command",
				Value: bson.M{
					"$regex": opt.StringValue(),
				},
			})
		}

		if opt, ok := m.Options["대답"]; ok {
			filter = append(filter, bson.E{
				Key: "result",
				Value: bson.M{
					"$regex": opt.StringValue(),
				},
			})
		}

		if opt, ok := m.Options["개수"]; ok {
			length = int(opt.IntValue())
		}
	}

	cur, err := databases.Database.Learns.Find(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ 오류",
				Description: "당신은 지식ㅇ을 가르쳐준 적이 없어요!",
				Color:       utils.EmbedFail,
			}

			switch m := m.(type) {
			case *discordgo.MessageCreate:
				s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
			case *utils.InteractionCreate:
				m.EditReply(&discordgo.WebhookEdit{
					Embeds: &[]*discordgo.MessageEmbed{embed},
				})
			}
			return
		}

		fmt.Println(err)
		embed := &discordgo.MessageEmbed{
			Title:       "❌ 오류",
			Description: "데이터를 가져오는데 실패했어요.",
			Color:       utils.EmbedFail,
		}

		switch m := m.(type) {
		case *discordgo.MessageCreate:
			s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
		case *utils.InteractionCreate:
			m.EditReply(&discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{embed},
			})
		}
		return
	}

	defer cur.Close(context.TODO())

	cur.All(context.TODO(), &data)

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s님이 알려주신 지식", globalName),
		Color: utils.EmbedDefault,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatarUrl,
		},
	}

	utils.StartPaginationEmbed(s, m, embed, getDescriptions(&data, length), utils.CodeBlock("md", fmt.Sprintf("# 총 %d개에요.\n", len(data))+"%s"))
}
