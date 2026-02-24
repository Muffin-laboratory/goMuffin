package chat

import (
	"log/slog"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/chatbot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

func Chat(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	var attachment discord.Attachment

	if opt, ok := data.OptAttachment("첨부파일"); ok {
		attachment = opt
	}

	content, err := chatbot.GetChatBot().GetResponse(e.Ctx, e.User(), data.String("내용"), attachment)
	if err != nil {
		slog.Error("failed to respond chat.", "user_id", e.User().ID, "error", err)
		_, err := e.UpdateInteractionResponse(
			discord.NewMessageUpdateV2([]discord.LayoutComponent{
				builders.MakeErrorContainer("대답하는 데 실패했어요. 잠시 후에 다시 시도해주세요."),
			}),
		)
		return err
	}

	result := chatbot.ParseResult(content, e)
	_, err = e.UpdateInteractionResponse(
		discord.NewMessageUpdate().
			WithContent(result).
			WithAllowedMentions(&discord.AllowedMentions{
				Parse: make([]discord.AllowedMentionType, 0),
				Users: make([]snowflake.ID, 0),
				Roles: make([]snowflake.ID, 0),
			}),
	)
	return err
}
