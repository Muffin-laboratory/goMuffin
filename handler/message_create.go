package handler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
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

// MessageCreate is handlers of messageCreate event
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	config := configs.GetConfig()
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, config.Bot.Prefix) {
		content := strings.TrimPrefix(m.Content, config.Bot.Prefix)
		args := argParser(content)
		command := commands.GetDiscommand().Aliases[args[0]]

		if command == "" {
			if !databases.GetDatabase().IsUser(m.Author.ID) {
				utils.NewMessageSender(&utils.MessageCreate{
					MessageCreate: m,
					Session:       s,
				}).
					AddComponents(utils.GetUserIsNotRegisteredErrContainer(configs.GetConfig().Bot.Prefix)).
					SetComponentsV2(true).
					SetReply(true).
					Send()
				return
			}

			blocked, reason := databases.GetDatabase().IsUserBlocked(m.Author.ID)
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

			str, err := chatbot.ChatBot.GetResponse(m.Author, strings.TrimPrefix(content, "대화 "))
			if err != nil {
				log.Println(err)
				utils.NewMessageSender(&utils.MessageCreate{
					MessageCreate: m,
					Session:       s,
				}).
					SetContent(str).
					SetReply(true).
					Send()
				return
			}

			result := chatbot.ParseResult(str, s, m)
			utils.NewMessageSender(&utils.MessageCreate{
				MessageCreate: m,
				Session:       s,
			}).
				SetContent(result).
				SetReply(true).
				Send()
			return
		}

		err := commands.GetDiscommand().MessageRun(command, s, m, args[1:])
		if err != nil {
			owner, _ := s.User(configs.GetConfig().Bot.OwnerID)
			log.Println(err)
			utils.NewMessageSender(&utils.MessageCreate{
				MessageCreate: m,
				Session:       s,
			}).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: fmt.Sprintf("오류가 발생하였어요. 만약 계속 발생한다면, %s으로 연락해주세요.", utils.InlineCode(owner.Username))})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
			return
		}
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
}
