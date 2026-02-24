package events

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/chatbot"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

// OnMessageCreate is handlers of messageCreate event
func OnMessageCreate(m *events.MessageCreate) {
	config := configs.Configs()
	if m.Message.Author.Bot {
		return
	}

	authorID := m.Message.Author.ID
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if content, ok := strings.CutPrefix(m.Message.Content, config.Bot.Prefix); ok {
		if !repository.GetDatabase().Users.IsUser(ctx, int64(authorID)) {
			if _, err := m.Client().Rest.CreateMessage(
				m.ChannelID,
				discord.NewMessageCreateV2(builders.MakeUserIsNotRegisteredErrContainer()).
					WithMessageReference(m.Message.MessageReference),
			); err != nil {
				slog.Error("failed to send you are not a user msg", "error", err)
			}

			return
		}

		if blocked, reason := repository.GetDatabase().Users.IsUserBlocked(ctx, int64(authorID)); blocked {
			if _, err := m.Client().Rest.CreateMessage(
				m.ChannelID,
				discord.NewMessageCreateV2(
					builders.MakeUserIsBlockedContainer(*m.Message.Author.GlobalName, reason),
				).
					WithMessageReference(m.Message.MessageReference),
			); err != nil {
				slog.Error("failed to send you are blocked msg", "error", err)
			}

			return
		}

		if err := m.Client().Rest.SendTyping(m.ChannelID); err != nil {
			slog.Error("failed to send typing", "error", err)
		}

		dbUser, err := repository.GetDatabase().Users.FindByID(ctx, int64(authorID))
		if err != nil {
			owner, _ := m.Client().Rest.GetUser(config.Bot.OwnerID)
			m.Client().Logger.Error("failed to get a user", "error", err)
			if _, err := m.Client().Rest.CreateMessage(
				m.ChannelID,
				discord.NewMessageCreateV2(
					builders.MakeErrorContainer("오류가 발생하였어요. 만약 계속 발생한다면, `%s`으로 연락해주세요.", owner.Username),
				).
					WithMessageReference(m.Message.MessageReference),
			); err != nil {
				slog.Error("failed to send err msg", "error", err)
			}

			return
		}

		str, err := chatbot.GetChatBot().GetResponse(ctx, m.Message.Author, content, m.Message.Attachments...)
		if err != nil {
			slog.Error("failed to respond chat.", "user_id", m.Message.Author.ID, "error", err)
			if _, err := m.Client().Rest.CreateMessage(
				m.ChannelID,
				discord.NewMessageCreateV2(
					builders.MakeErrorContainer("대답하는 데 실패했어요. 잠시 후에 다시 시도해주세요."),
				).
					WithMessageReference(m.Message.MessageReference).
					WithAllowedMentions(&discord.AllowedMentions{
						Parse:       make([]discord.AllowedMentionType, 0),
						Users:       make([]snowflake.ID, 0),
						Roles:       make([]snowflake.ID, 0),
						RepliedUser: dbUser.ReplyUser,
					}),
			); err != nil {
				slog.Error("failed to send err msg", "error", err)
			}

			return
		}

		result := chatbot.ParseResult(str, m)
		if _, err = m.Client().Rest.CreateMessage(
			m.ChannelID,
			discord.NewMessageCreate().
				WithContent(result).
				WithMessageReference(m.Message.MessageReference).
				WithAllowedMentions(&discord.AllowedMentions{
					Parse:       make([]discord.AllowedMentionType, 0),
					Users:       make([]snowflake.ID, 0),
					Roles:       make([]snowflake.ID, 0),
					RepliedUser: dbUser.ReplyUser,
				}),
		); err != nil {
			slog.Error("failed to send chat msg", "error", err)
		}

		return
	} else {
		if m.Message.Author.ID == config.Chatbot.Train.UserID {
			if _, err := repository.GetDatabase().Texts.Create(ctx, m.Message.Content); err != nil {
				slog.Error("failed to save muffin data.", "error", err)
			}
		}
	}
}
