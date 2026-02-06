package handler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/chatbot"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/bwmarrin/discordgo"
)

// MessageCreate is handlers of messageCreate event
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	config := configs.GetConfig()
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if strings.HasPrefix(m.Content, config.Bot.Prefix) {

		m := &builders.MessageCreate{
			MessageCreate: m,
			Session:       s,
			Ctx:           ctx,
		}
		content := strings.TrimPrefix(m.Content, config.Bot.Prefix)

		if !repository.GetDatabase().Users.IsUser(m.Ctx, m.Author.ID) {
			builders.NewMessageSender(m).
				AddComponents(builders.MakeUserIsNotRegisteredErrContainer()).
				SetComponentsV2(true).
				SetReply(true).
				Send()

			return
		}

		blocked, reason := repository.GetDatabase().Users.IsUserBlocked(m.Ctx, m.Author.ID)
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

		dbUser, err := repository.GetDatabase().Users.FindByID(m.Ctx, m.Author.ID)
		if err != nil {
			owner, _ := s.User(configs.GetConfig().Bot.OwnerID)
			log.Println(err)
			builders.NewMessageSender(m).
				AddComponents(builders.MakeErrorContainer(fmt.Sprintf("오류가 발생하였어요. 만약 계속 발생한다면, %s으로 연락해주세요.", utils.InlineCode(owner.Username)))).
				SetComponentsV2(true).
				SetReply(true).
				Send()

			return
		}

		str, err := chatbot.GetChatBot().GetResponse(m.Ctx, m.Author, content, m.Attachments...)
		if err != nil {
			log.Println(err)
			builders.NewMessageSender(m).
				SetContent(str).
				SetReply(true).
				SetAllowedMentions(discordgo.MessageAllowedMentions{
					Parse:       []discordgo.AllowedMentionType{},
					Users:       []string{},
					Roles:       []string{},
					RepliedUser: dbUser.ReplyUser,
				}).
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
				RepliedUser: dbUser.ReplyUser,
			}).
			Send()
	} else {
		if m.Author.ID == config.Chatbot.Train.UserID {
			if _, err := repository.GetDatabase().Texts.Create(ctx, m.Content); err != nil {
				log.Println(err)
			}
		}
	}
}
