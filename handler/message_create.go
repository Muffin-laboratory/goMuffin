package handler

import (
	"context"
	"log"
	"strings"
	"time"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
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
			m := &utils.MessageCreate{
				MessageCreate: m,
				Session:       s,
			}
			content := strings.TrimPrefix(m.Content, config.Bot.Prefix)

			if !databases.GetDatabase().Users.IsUser(m.Author.ID) {
				utils.NewMessageSender(m).
					AddComponents(utils.GetUserIsNotRegisteredErrContainer(config.Bot.Prefix)).
					SetComponentsV2(true).
					SetReply(true).
					Send()
				return
			}

			blocked, reason := databases.GetDatabase().Users.IsUserBlocked(m.Author.ID)
			if blocked {
				user, _ := s.User(m.Author.ID)
				utils.NewMessageSender(m).
					AddComponents(utils.GetUserIsBlockedContainer(user.GlobalName, reason)).
					SetComponentsV2(true).
					SetReply(true).
					Send()
				return
			}

			s.ChannelTyping(m.ChannelID)

			str, err := chatbot.GetChatBot().GetResponse(m.Author, content)
			if err != nil {
				log.Println(err)
				utils.NewMessageSender(m).
					SetContent(str).
					SetReply(true).
					Send()

				return
			}

			result := chatbot.ParseResult(str, s, m)
			utils.NewMessageSender(m).
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
