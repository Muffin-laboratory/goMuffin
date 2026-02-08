package chat

import (
	"context"
	"log/slog"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/chatbot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func Chat(ctx context.Context, i *builders.CommandCreate) error {
	var attachment discord.Attachment

	commandData := i.SlashCommandInteractionData()

	if opt, ok := commandData.OptAttachment("첨부파일"); ok {
		attachment = opt
	}

	content, err := chatbot.GetChatBot().GetResponse(ctx, i.User(), commandData.String("내용"), attachment)
	if err != nil {
		slog.Error("error in responding chat.", "user_id", i.User().ID, "error", err)
		builders.NewMessageSender(i).
			SetContent(content).
			Send()
		return nil
	}

	result := chatbot.ParseResult(content, i)
	return builders.NewMessageSender(i).
		SetContent(result).
		SetAllowedMentions(discord.AllowedMentions{
			Parse: make([]discord.AllowedMentionType, 0),
			Users: make([]snowflake.ID, 0),
			Roles: make([]snowflake.ID, 0),
		}).
		Send()
}
