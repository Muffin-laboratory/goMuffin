package knowledge

import (
	"context"
	"fmt"
	"strings"

	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
)

func Learn(ctx context.Context, i *builders.CommandCreate, igCommands []string) error {
	command := i.SlashCommandInteractionData().String("단어")
	result := i.SlashCommandInteractionData().String("대답")

	ignores := []string{"미간", "Migan", "migan", "간미"}
	ignores = append(ignores, igCommands...)

	disallows := []string{
		"@everyone",
		"@here",
		fmt.Sprintf("<@%s>", configs.GetConfig().Bot.OwnerID),
	}

	for _, ig := range ignores {
		if strings.Contains(command, ig) {
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeErrorContainer("해당 단어는 배우기 껄끄럽네요.")).
				SetComponentsV2(true).
				SetReply(true).
				Send()
		}
	}

	for _, di := range disallows {
		if strings.Contains(result, di) {
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeErrorContainer("해당 단어의 대답으로 하기 좀 그렇네요.")).
				SetComponentsV2(true).
				SetReply(true).
				Send()
		}
	}

	if _, err := repository.GetDatabase().Knowledge.Create(ctx, i.User().ID.String(), command, result); err != nil {
		return err
	}

	return builders.NewMessageSender(i).
		AddComponents(builders.MakeSuccessContainer("%s 배웠어요.", hangul.GetJosa(command, hangul.EUL_REUL))).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
