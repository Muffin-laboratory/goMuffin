package commands

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type chStruct struct {
	name   dataType
	length int
}

type dataType int

const (
	text dataType = iota
	muffin
	nsfw
	learn
	userLearn
)

var dataLengthWg sync.WaitGroup

var DataLengthCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "데이터학습량",
		Description: "봇이 학습한 데ㅇ이터량을 보여줘요.",
	},
	Aliases: []string{"학습데이터량", "데이터량", "학습량"},
	DetailedDescription: &DetailedDescription{
		Usage: fmt.Sprintf("%s학습데이터량", configs.Config.Bot.Prefix),
	},
	Category:                   General,
	RegisterApplicationCommand: true,
	RegisterMessageCommand:     true,
	Flags:                      CommandFlagsIsRegistered,
	MessageRun: func(ctx *MsgContext) error {
		return dataLengthRun(ctx.Msg.Session, ctx.Msg, ctx.Msg.Author.Username, ctx.Msg.Author.ID)
	},
	ChatInputRun: func(ctx *ChatInputContext) error {
		ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		})

		return dataLengthRun(ctx.Inter.Session, ctx.Inter, ctx.Inter.Member.User.Username, ctx.Inter.Member.User.ID)
	},
}

func getLength(ch chan chStruct, dType dataType, coll *mongo.Collection, filter bson.D) {
	defer dataLengthWg.Done()
	var err error
	var cur *mongo.Cursor
	var data []bson.M

	cur, err = coll.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}

	defer cur.Close(context.TODO())

	cur.All(context.TODO(), &data)
	ch <- chStruct{name: dType, length: len(data)}
}

func dataLengthRun(s *discordgo.Session, m any, username, userId string) error {
	ch := make(chan chStruct)
	var textLength,
		muffinLength,
		nsfwLength,
		learnLength,
		userLearnLength int

	dataLengthWg.Add(5)
	go getLength(ch, text, databases.Database.Texts, bson.D{{}})
	go getLength(ch, muffin, databases.Database.Texts, bson.D{{Key: "persona", Value: "muffin"}})
	go getLength(ch, nsfw, databases.Database.Texts, bson.D{
		{
			Key: "persona",
			Value: bson.M{
				"$regex": "^user",
			},
		},
	})
	go getLength(ch, learn, databases.Database.Learns, bson.D{{}})
	go getLength(ch, userLearn, databases.Database.Learns, bson.D{{Key: "user_id", Value: userId}})

	go func() {
		dataLengthWg.Wait()
		close(ch)
	}()

	for resp := range ch {
		switch dataType(resp.name) {
		case text:
			textLength = resp.length
		case muffin:
			muffinLength = resp.length
		case nsfw:
			nsfwLength = resp.length
		case learn:
			learnLength = resp.length
		case userLearn:
			userLearnLength = resp.length
		}
	}

	sum := textLength + learnLength

	return utils.NewMessageSender(m).
		AddComponents(discordgo.Container{
			Components: []discordgo.MessageComponent{
				discordgo.Section{
					Accessory: discordgo.Thumbnail{
						Media: discordgo.UnfurledMediaItem{
							URL: s.State.User.AvatarURL("512"),
						},
					},
					Components: []discordgo.MessageComponent{
						discordgo.TextDisplay{
							Content: fmt.Sprintf("### 저장된 데이터량\n총합: %s", utils.InlineCode(strconv.Itoa(sum)+"개")),
						},
						discordgo.TextDisplay{
							Content: fmt.Sprintf("- **총 채팅 데이터량**\n> %s", utils.InlineCode(strconv.Itoa(textLength))+"개"),
						},
						discordgo.TextDisplay{
							Content: fmt.Sprintf("- **머핀 데이터량**\n> %s", utils.InlineCode(strconv.Itoa(muffinLength))+"개"),
						},
					},
				},
				discordgo.TextDisplay{
					Content: fmt.Sprintf("- **nsfw 데이터량**\n> %s", utils.InlineCode(strconv.Itoa(nsfwLength))+"개"),
				},
				discordgo.TextDisplay{
					Content: fmt.Sprintf("- **총 지식 데이터량**\n> %s", utils.InlineCode(strconv.Itoa(learnLength))+"개"),
				},
				discordgo.TextDisplay{
					Content: fmt.Sprintf("- **%s님이 가르쳐준 데이터량**\n> %s", username, utils.InlineCode(strconv.Itoa(userLearnLength))+"개"),
				},
			},
		}).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
