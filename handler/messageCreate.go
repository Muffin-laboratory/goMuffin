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

		if command == "" {
			s.ChannelTyping(m.ChannelID)

			result := chatbot.ParseResult(chatbot.ChatBot.GetResponse(content), s, m)
			err := utils.NewMessageSender(m).
				SetContent(result).
				SetReply(true).
				Send()
			fmt.Println(err)
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
