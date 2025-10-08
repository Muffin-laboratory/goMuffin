package handler

import (
	"context"
	"log"
	"strings"
	"time"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"github.com/bwmarrin/discordgo"
)

// MessageCreate is handlers of messageCreate event
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	go func() {
		config := configs.GetConfig()
		if m.Author.ID == s.State.User.ID || m.Author.Bot {
			return
		}

		if strings.HasPrefix(m.Content, config.Bot.Prefix) {
			m := &builders.MessageCreate{
				MessageCreate: m,
				Session:       s,
			}
			content := strings.TrimPrefix(m.Content, config.Bot.Prefix)

			if !databases.GetDatabase().Users.IsUser(m.Author.ID) {
				builders.NewMessageSender(m).
					AddComponents(builders.MakeUserIsNotRegisteredErrContainer()).
					SetComponentsV2(true).
					SetReply(true).
					Send()
				return
			}

			blocked, reason := databases.GetDatabase().Users.IsUserBlocked(m.Author.ID)
			if blocked {
				user, _ := s.User(m.Author.ID)
				builders.NewMessageSender(m).
					AddComponents(builders.MakeUserIsBlockedContainer(user.GlobalName, reason)).
					SetComponentsV2(true).
					SetReply(true).
					Send()
				return
			}

			s.ChannelTyping(m.ChannelID)

			str, err := chatbot.GetChatBot().GetResponse(m.Author, content)
			if err != nil {
				log.Println(err)
				builders.NewMessageSender(m).
					SetContent(str).
					SetReply(true).
					Send()

				return
			}

			result := chatbot.ParseResult(str, s, m)
			builders.NewMessageSender(m).
				SetContent(result).
				SetReply(true).
				SetAllowedMentions(discordgo.MessageAllowedMentions{
					Parse:       []discordgo.AllowedMentionType{},
					Users:       []string{},
					Roles:       []string{},
					RepliedUser: true,
				}).
				Send()

			return
		} else {
			if m.Author.ID == config.Chatbot.Train.UserID {
				if _, err := databases.GetDatabase().Texts.InsertOne(context.TODO(), databases.Text{
					Text:      m.Content,
					Persona:   "muffin",
					CreatedAt: time.Now(),
				}); err != nil {
					log.Fatalln(err)
				}
			}

			return
		}
	}()
}
