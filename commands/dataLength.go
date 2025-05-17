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

// var dataLengthCh chan chStruct = make(chan chStruct)
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
	Category: General,
	MessageRun: func(ctx *MsgContext) {
		dataLengthRun(ctx.Msg, ctx.Msg.Author.Username, ctx.Msg.Author.ID)
	},
	ChatInputRun: func(ctx *ChatInputContext) {
		ctx.Inter.DeferReply(true)
		dataLengthRun(ctx.Inter, ctx.Inter.Member.User.Username, ctx.Inter.Member.User.ID)
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

func dataLengthRun(m any, username, userId string) {
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

	// 나중에 djs처럼 Embed 만들어 주는 함수 만들어야겠다
	// 지금은 임시방편
	utils.NewMessageSender(m).
		AddEmbeds(&discordgo.MessageEmbed{
			Title:       "저장된 데이터량",
			Description: fmt.Sprintf("총합: %s개", utils.InlineCode(strconv.Itoa(sum))),
			Color:       utils.EmbedDefault,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "총 채팅 데이터량",
					Value:  utils.InlineCode(strconv.Itoa(textLength)) + "개",
					Inline: true,
				},
				{
					Name:   "총 지식 데이터량",
					Value:  utils.InlineCode(strconv.Itoa(learnLength)) + "개",
					Inline: true,
				},
				{
					Name:  "머핀 데이터량",
					Value: utils.InlineCode(strconv.Itoa(muffinLength)) + "개",
				},
				{
					Name:   "nsfw 데이터량",
					Value:  utils.InlineCode(strconv.Itoa(nsfwLength)) + "개",
					Inline: true,
				},
				{
					Name:   fmt.Sprintf("%s님이 가르쳐준 데이터량", username),
					Value:  utils.InlineCode(strconv.Itoa(userLearnLength)) + "개",
					Inline: true,
				},
			},
		}).
		SetReply(true).
		Send()
}
