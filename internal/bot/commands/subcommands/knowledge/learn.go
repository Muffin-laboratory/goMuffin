package knowledge

import (
	"fmt"
	"strings"

	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func Learn(i *builders.InteractionCreate, opts *discordgo.ApplicationCommandInteractionDataOption, igCommands []string) error {
	command := opts.GetOption("단어").StringValue()
	result := opts.GetOption("대답").StringValue()

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

	if _, err := repository.GetDatabase().Knowledge.Create(i.Ctx, i.User.ID, command, result); err != nil {
		return err
	}

	return builders.NewMessageSender(i).
		AddComponents(builders.MakeSuccessContainer("%s 배웠어요.", hangul.GetJosa(command, hangul.EUL_REUL))).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
