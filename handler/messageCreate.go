package handler

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func argParser(content string) (args []string) {
	for _, arg := range utils.RegexpFlexibleString.FindAllStringSubmatch(content, -1) {
		if arg[1] != "" {
			args = append(args, arg[1])
		} else {
			args = append(args, arg[0])
		}
	}
	return
}

func resultParser(content string, s *discordgo.Session, m *discordgo.MessageCreate) string {
	result := content
	userCreatedAt, _ := discordgo.SnowflakeTimestamp(m.Author.ID)

	result = strings.ReplaceAll(result, "{user.name}", m.Author.Username)
	result = strings.ReplaceAll(result, "{user.mention}", m.Author.Mention())
	result = strings.ReplaceAll(result, "{user.globalName}", m.Author.GlobalName)
	result = strings.ReplaceAll(result, "{user.id}", m.Author.ID)
	result = strings.ReplaceAll(result, "{user.createdAt}", utils.Time(&userCreatedAt, utils.RelativeTime))
	result = strings.ReplaceAll(result, "{user.joinedAt}", utils.Time(&m.Member.JoinedAt, utils.RelativeTime))

	result = strings.ReplaceAll(result, "{muffin.version}", configs.MUFFIN_VERSION)
	result = strings.ReplaceAll(result, "{muffin.updatedAt}", utils.Time(configs.UpdatedAt, utils.RelativeTime))
	result = strings.ReplaceAll(result, "{muffin.startedAt}", utils.Time(configs.StartedAt, utils.RelativeTime))
	result = strings.ReplaceAll(result, "{muffin.name}", s.State.User.Username)
	result = strings.ReplaceAll(result, "{muffin.id}", s.State.User.ID)
	return result
}

// MessageCreate is handlers of messageCreate event
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	config := configs.Config
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, config.Bot.Prefix) {
		content := strings.TrimPrefix(m.Content, config.Bot.Prefix)
		args := argParser(content)
		command := commands.Discommand.Aliases[args[0]]

		if m.Author.ID == config.Train.UserID {
			if _, err := databases.Database.Texts.InsertOne(context.TODO(), databases.InsertText{
				Text:      content,
				Persona:   "muffin",
				CreatedAt: time.Now(),
			}); err != nil {
				log.Fatalln(err)
			}
		}

		if command == "" {
			s.ChannelTyping(m.ChannelID)

			var data []databases.Text
			var learnData []databases.Learn
			var filter bson.D

			ch := make(chan int)
			x := rand.Intn(10)

			channel, _ := s.Channel(m.ChannelID)
			if channel.NSFW {
				filter = bson.D{{}}

				if _, err := databases.Database.Texts.InsertOne(context.TODO(), databases.InsertText{
					Text:      content,
					Persona:   fmt.Sprintf("user:%s", m.Author.Username),
					CreatedAt: time.Now(),
				}); err != nil {
					log.Fatalln(err)
				}

			} else {
				filter = bson.D{{Key: "persona", Value: "muffin"}}
			}

			go func() {
				cur, err := databases.Database.Texts.Find(context.TODO(), filter)
				if err != nil {
					log.Fatalln(err)
				}

				defer cur.Close(context.TODO())

				cur.All(context.TODO(), &data)
				ch <- 1
			}()
			go func() {
				cur, err := databases.Database.Learns.Find(context.TODO(), bson.D{{Key: "command", Value: content}})
				if err != nil {
					if err == mongo.ErrNilDocument {
						learnData = []databases.Learn{}
					}
					log.Fatalln(err)
				}

				defer cur.Close(context.TODO())

				cur.All(context.TODO(), &learnData)
				ch <- 1
			}()

			for range 2 {
				<-ch
			}
			close(ch)

			if x > 2 && len(learnData) != 0 {
				data := learnData[rand.Intn(len(learnData))]
				user, _ := s.User(data.UserId)
				result := resultParser(data.Result, s, m)

				s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
					Reference: m.Reference(),
					Content:   fmt.Sprintf("%s\n%s", result, utils.InlineCode(fmt.Sprintf("%s님이 알려주셨어요.", user.Username))),
					AllowedMentions: &discordgo.MessageAllowedMentions{
						Roles: []string{},
						Parse: []discordgo.AllowedMentionType{},
						Users: []string{},
					},
				})
				return
			}

			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Reference: m.Reference(),
				Content:   data[rand.Intn(len(data))].Text,
				AllowedMentions: &discordgo.MessageAllowedMentions{
					Roles: []string{},
					Parse: []discordgo.AllowedMentionType{},
					Users: []string{},
				},
			})
			return
		}

		commands.Discommand.MessageRun(command, s, m, args[1:])
		return
	} else {
		return
	}
}
