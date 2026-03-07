package commands

import (
	"fmt"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	const name = "핑"

	loader.GetDiscommand().RegisterCommand(discord.SlashCommandCreate{
		Name:        name,
		Description: "봇의 레이턴시를 확인해요.",
	})

	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Command("/"+name, func(e *handler.CommandEvent) error {
			var dbPing int64

			bot, _ := e.Client().Caches.SelfUser()
			title := fmt.Sprintf("### 🏓 %s의 지연시간", bot.Username)

			start := time.Now()
			if err := repository.GetDatabase().Client.Ping(e.Ctx, nil); err != nil {
				return err
			}
			dbPing = time.Since(start).Milliseconds()

			if err := e.CreateMessage(
				discord.NewMessageCreateV2(
					discord.NewContainer(
						discord.NewTextDisplay(title),
						discord.NewTextDisplay("- 지연시간 측정 중..."),
					),
				),
			); err != nil {
				return err
			}

			message, err := e.Client().Rest.GetInteractionResponse(e.ApplicationID(), e.Token())
			if err != nil {
				return err
			}

			createdTimestamp := e.ID().Time()
			discordPing := message.ID.Time().Sub(createdTimestamp).Milliseconds()

			_, err = e.UpdateInteractionResponse(
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					discord.NewContainer(
						discord.NewTextDisplay(title),
						discord.NewTextDisplayf("- **디스코드 지연시간:** `%d`ms", discordPing),
						discord.NewTextDisplayf("- **데이터베이스 지연시간:** `%d`ms", dbPing),
					),
				}),
			)
			return err
		})
	})
}
