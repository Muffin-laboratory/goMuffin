package knowledge

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/LoperLee/golang-hangul-toolkit/hangul"
	"github.com/bwmarrin/discordgo"
)

func Learn(i *utils.InteractionCreate, opts utils.CommandInteractionOptionsMap, igCommands []string) error {
	err := i.DeferReply(&discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		return err
	}

	command := opts["단어"].StringValue()
	result := opts["대답"].StringValue()

	ignores := []string{"미간", "Migan", "migan", "간미"}
	ignores = append(ignores, igCommands...)

	disallows := []string{
		"@everyone",
		"@here",
		fmt.Sprintf("<@%s>", configs.GetConfig().Bot.OwnerID),
	}

	for _, ig := range ignores {
		if strings.Contains(command, ig) {
			return utils.NewMessageSender(i).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "해ㄷ당 단어는 배우기 껄끄럽네요."})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
		}
	}

	for _, di := range disallows {
		if strings.Contains(result, di) {
			return utils.NewMessageSender(i).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "해당 단ㅇ어의 대답으로 하기 좀 그렇네요."})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
		}
	}

	_, err = databases.GetDatabase().Knowledge.InsertOne(context.TODO(), databases.Knowledge{
		Command:   command,
		Result:    result,
		UserID:    i.User.ID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return utils.NewMessageSender(i).
		AddComponents(utils.GetSuccessContainer(
			discordgo.TextDisplay{
				Content: fmt.Sprintf("%s 배웠어요.", hangul.GetJosa(command, hangul.EUL_REUL)),
			},
		)).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
