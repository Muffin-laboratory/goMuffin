package commands

import (
	"context"
	"fmt"
	"strings"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var LearnedDataListCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "리스트",
		Description: "당신이 가ㄹ르쳐준 단어를 나열해요.",
	},
	Aliases: []string{"list", "목록", "지식목록"},
	DetailedDescription: &DetailedDescription{
		Usage: "머핀아 리스트",
	},
	Category: Chattings,
	MessageRun: func(ctx *MsgContext) {
		learnedDataListRun(ctx.Session, ctx.Msg)
	},
	ChatInputRun: func(ctx *ChatInputContext) {
		learnedDataListRun(ctx.Session, ctx.Inter)
	},
}

func getDescriptions(datas *[]databases.Learn) (descriptions []string) {
	for _, data := range *datas {
		descriptions = append(descriptions, fmt.Sprintf("- %s: %s", data.Command, data.Result))
	}
	return
}

func learnedDataListRun(s *discordgo.Session, m any) {
	var userId, globalName, avatarUrl string
	var datas []databases.Learn
	switch m := m.(type) {
	case *discordgo.MessageCreate:
		userId = m.Author.ID
		globalName = m.Author.GlobalName
		avatarUrl = m.Author.AvatarURL("512")
	case *discordgo.InteractionCreate:
		s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})

		userId = m.Member.User.ID
		globalName = m.Member.User.GlobalName
		avatarUrl = m.User.AvatarURL("512")
	}

	cur, err := databases.Learns.Find(context.TODO(), bson.D{{Key: "user_id", Value: userId}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ 오류",
				Description: "당신은 지식ㅇ을 가르쳐준 적이 없어요!",
				Color:       int(utils.EFail),
			}

			switch m := m.(type) {
			case *discordgo.MessageCreate:
				s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
			case *discordgo.InteractionCreate:
				s.InteractionResponseEdit(m.Interaction, &discordgo.WebhookEdit{
					Embeds: &[]*discordgo.MessageEmbed{embed},
				})
			}
			return
		}

		fmt.Println(err)
		embed := &discordgo.MessageEmbed{
			Title:       "❌ 오류",
			Description: "데이터를 가져오는데 실패했어요.",
			Color:       int(utils.EFail),
		}

		switch m := m.(type) {
		case *discordgo.MessageCreate:
			s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
		case *discordgo.InteractionCreate:
			s.InteractionResponseEdit(m.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{embed},
			})
		}
		return
	}

	defer cur.Close(context.TODO())

	cur.All(context.TODO(), &datas)

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s님이 알려주신 지식", globalName),
		Description: utils.CodeBlockWithLanguage("md", fmt.Sprintf("# 총 %d개에요.\n%s", len(datas), strings.Join(getDescriptions(&datas), "\n"))),
		Color:       int(utils.EDefault),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatarUrl,
		},
	}

	switch m := m.(type) {
	case *discordgo.MessageCreate:
		s.ChannelMessageSendEmbedReply(m.ChannelID, embed, m.Reference())
	case *discordgo.InteractionCreate:
		s.InteractionResponseEdit(m.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		})
	}
}
