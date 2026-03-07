package knowledge

import (
	"fmt"
	"strings"

	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func Learn(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	command := data.String("단어")
	result := data.String("대답")

	ignores := []string{"미간", "Migan", "migan", "간미"}
	disallows := []string{
		"@everyone",
		"@here",
		fmt.Sprintf("<@%s>", configs.Configs().Bot.OwnerID),
	}

	for _, ig := range ignores {
		if strings.Contains(command, ig) {
			_, err := e.UpdateInteractionResponse(
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					builders.MakeErrorContainer("해당 단어는 배우기 껄끄럽네요."),
				}),
			)
			return err
		}
	}

	for _, di := range disallows {
		if strings.Contains(result, di) {
			_, err := e.UpdateInteractionResponse(
				discord.NewMessageUpdateV2([]discord.LayoutComponent{
					builders.MakeErrorContainer("해당 단어의 대답으로 하기 좀 그렇네요."),
				}),
			)
			return err
		}
	}

	if _, err := repository.GetDatabase().Knowledge.Create(e.Ctx, int64(e.User().ID), command, result); err != nil {
		return err
	}

	_, err := e.UpdateInteractionResponse(
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			builders.MakeSuccessContainer("%s 배웠어요.", hangul.GetJosa(command, hangul.EUL_REUL)),
		}),
	)
	return err
}
