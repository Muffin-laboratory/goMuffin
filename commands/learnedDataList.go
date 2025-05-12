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
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	learnArgsCommand = "단어:"
	learnArgsResult  = "대답:"
)

var LearnedDataListCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "리스트",
		Description: "당신이 가ㄹ르쳐준 지식을 나열해요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "쿼리",
				Description: "해당 단어가 포함된 결과만 찾아요.",
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

func getDescriptions(data *[]databases.Learn) (descriptions []string) {
	MAX_LENGTH := 100

	for _, data := range *data {
		command := data.Command
		result := data.Result

		if runeCommand := []rune(command); len(runeCommand) >= MAX_LENGTH {
			command = string(runeCommand)[:MAX_LENGTH] + "..."
		}

		if runeResult := []rune(result); len(runeResult) >= MAX_LENGTH {
			result = string(runeResult[:MAX_LENGTH]) + "..."
		}

		descriptions = append(descriptions, fmt.Sprintf("- %s: %s", command, result))
	}
	return
}

func learnedDataListRun(s *discordgo.Session, m any, args *[]string) {
	var userId, globalName, avatarUrl string
	var data []databases.Learn
	var filter bson.E

	switch m := m.(type) {
	case *discordgo.MessageCreate:
		userId = m.Author.ID
		globalName = m.Author.GlobalName
		avatarUrl = m.Author.AvatarURL("512")

		query := strings.Join(*args, " ")
		if strings.HasPrefix(query, learnArgsResult) {
			query, _ = strings.CutPrefix(query, learnArgsResult)
			filter = bson.E{
				Key: "result",
				Value: bson.M{
					"$regex": query,
				},
			}
		} else {
			query, _ = strings.CutPrefix(query, learnArgsCommand)
			filter = bson.E{
				Key: "command",
				Value: bson.M{
					"$regex": query,
				},
			}
		}
	case *utils.InteractionCreate:
		m.DeferReply(true)

		userId = m.Member.User.ID
		globalName = m.Member.User.GlobalName
		avatarUrl = m.Member.User.AvatarURL("512")

		if opt, ok := m.Options["쿼리"]; ok {
			query := opt.StringValue()

			if strings.HasPrefix(query, learnArgsResult) {
				query, _ = strings.CutPrefix(query, learnArgsResult)
				filter = bson.E{
					Key: "result",
					Value: bson.M{
						"$regex": query,
					},
				}
			} else {
				query, _ = strings.CutPrefix(query, learnArgsCommand)
				filter = bson.E{
					Key: "command",
					Value: bson.M{
						"$regex": query,
					},
				}
			}
		}
	}

	cur, err := databases.Database.Learns.Find(context.TODO(), bson.D{{Key: "user_id", Value: userId}, filter})
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
		Title:       fmt.Sprintf("%s님이 알려주신 지식", globalName),
		Description: utils.CodeBlock("md", fmt.Sprintf("# 총 %d개에요.\n%s", len(data), strings.Join(getDescriptions(&data), "\n"))),
		Color:       utils.EmbedDefault,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatarUrl,
		},
	}

	switch m := m.(type) {
	case *discordgo.MessageCreate:
		s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
	case *utils.InteractionCreate:
		m.EditReply(&discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		})
	}
}
