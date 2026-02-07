package handler

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/chatbot"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

// OnMessageCreate is handlers of messageCreate event
func OnMessageCreate(m *events.MessageCreate) {
	config := configs.GetConfig()
	if m.Message.Author.Bot {
		return
	}

	authorID := m.Message.Author.ID

	if strings.HasPrefix(m.Message.Content, config.Bot.Prefix) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		content := strings.TrimPrefix(m.Message.Content, config.Bot.Prefix)

		if !repository.GetDatabase().Users.IsUser(ctx, authorID.String()) {
			m.Client().Rest.CreateMessage(
				m.ChannelID,
				discord.NewMessageCreateBuilder().
					SetComponents(builders.MakeUserIsNotRegisteredErrContainer()).
					SetIsComponentsV2(true).
					SetMessageReference(m.Message.MessageReference).
					Build(),
			)

			return
		}

		if blocked, reason := repository.GetDatabase().Users.IsUserBlocked(ctx, authorID.String()); blocked {
			m.Client().Rest.CreateMessage(
				m.ChannelID,
				discord.NewMessageCreateBuilder().
					SetComponents(builders.MakeUserIsBlockedContainer(*m.Message.Author.GlobalName, reason)).
					SetIsComponentsV2(true).
					SetMessageReference(m.Message.MessageReference).
					Build(),
			)

			return
		}

		m.Client().Rest.SendTyping(m.ChannelID)

		dbUser, err := repository.GetDatabase().Users.FindByID(ctx, authorID.String())
		if err != nil {
			owner, _ := m.Client().Rest.GetUser(snowflake.MustParse(config.Bot.OwnerID))
			m.Client().Logger.Error("%v", err)
			m.Client().Rest.CreateMessage(
				m.ChannelID,
				discord.NewMessageCreateBuilder().
					SetComponents(builders.MakeErrorContainer("오류가 발생하였어요. 만약 계속 발생한다면, %s으로 연락해주세요.", utils.InlineCode(owner.Username))).
					SetIsComponentsV2(true).
					SetMessageReference(m.Message.MessageReference).
					Build(),
			)

			return
		}

		str, err := chatbot.GetChatBot().GetResponse(ctx, m.Message.Author, content, m.Message.Attachments...)
		if err != nil {
			m.Client().Logger.Error("%v", err)
			m.Client().Rest.CreateMessage(
				m.ChannelID,
				discord.NewMessageCreateBuilder().
					SetContent(str).
					SetMessageReference(m.Message.MessageReference).
					SetAllowedMentions(&discord.AllowedMentions{
						Parse:       make([]discord.AllowedMentionType, 0),
						Users:       make([]snowflake.ID, 0),
						Roles:       make([]snowflake.ID, 0),
						RepliedUser: dbUser.ReplyUser,
					}).
					Build(),
			)

			return
		}

		result := chatbot.ParseResult(str, m)
		m.Client().Rest.CreateMessage(
			m.ChannelID,
			discord.NewMessageCreateBuilder().
				SetContent(result).
				SetMessageReference(m.Message.MessageReference).
				SetAllowedMentions(&discord.AllowedMentions{
					Parse:       make([]discord.AllowedMentionType, 0),
					Users:       make([]snowflake.ID, 0),
					Roles:       make([]snowflake.ID, 0),
					RepliedUser: dbUser.ReplyUser,
				}).
				Build(),
		)

		return
	} else {
		if m.Message.Author.ID.String() == config.Chatbot.Train.UserID {
			if _, err := repository.GetDatabase().Texts.InsertOne(context.TODO(), repository.Text{
				Text:      m.Message.Content,
				Persona:   "muffin",
				CreatedAt: time.Now(),
			}); err != nil {
				log.Fatalln(err)
			}
		}

		return
	}
}
