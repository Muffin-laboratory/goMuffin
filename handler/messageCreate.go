package handler

import (
	"context"
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
	config := configs.Config
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, config.Bot.Prefix) {
		content := strings.TrimPrefix(m.Content, config.Bot.Prefix)
		args := argParser(content)
		command := commands.Discommand.Aliases[args[0]]

		if command == "" {
			s.ChannelTyping(m.ChannelID)

			result := chatbot.ParseResult(chatbot.ChatBot.GetResponse(content), s, m)
			utils.NewMessageSender(&utils.MessageCreate{
				MessageCreate: m,
				Session:       s,
			}).
				SetContent(result).
				SetReply(true).
				Send()
			return
		}

		commands.Discommand.MessageRun(command, s, m, args[1:])
		return
	} else {
		if m.Author.ID == config.Train.UserID {
			if _, err := databases.Database.Texts.InsertOne(context.TODO(), databases.InsertText{
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
