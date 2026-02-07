package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
)

var PingCommand = &loader.Command{
	SlashCommandCreate: &discord.SlashCommandCreate{
		Name:        "핑",
		Description: "봇의 레이턴시를 확인해요.",
	},
	Run: func(ctx context.Context, inter *builders.CommandCreate) error {
		var dbPing int64

		bot, _ := inter.Client().Caches.SelfUser()
		title := fmt.Sprintf("### 🏓 %s의 지연시간", bot.Username)

		start := time.Now()
		if err := repository.GetDatabase().Client.Ping(ctx, nil); err != nil {
			return err
		}
		dbPing = time.Since(start).Milliseconds()

		if err := builders.NewMessageSender(inter).
			AddComponents(
				discord.NewContainer(
					discord.NewTextDisplay(title),
					discord.NewTextDisplay("- 지연시간 측정 중..."),
				),
			).
			SetComponentsV2(true).
			Send(); err != nil {
			return err
		}

		message, err := inter.Client().Rest.GetInteractionResponse(inter.ApplicationID(), inter.Token())
		if err != nil {
			return err
		}

		createdTimestamp := inter.ID().Time()
		discordPing := message.ID.Time().Sub(createdTimestamp)

		return builders.NewMessageSender(inter).
			AddComponents(
				discord.NewContainer(
					discord.NewTextDisplay(title),
					discord.NewTextDisplayf("- **디스코드 지연시간:** `%d`ms", discordPing),
					discord.NewTextDisplayf("- **데이터베이스 지연시간:** `%d`ms", dbPing),
				),
			).
			SetComponentsV2(true).
			Send()
	},
}

func init() {
	loader.GetDiscommand().LoadCommand(PingCommand)
}
