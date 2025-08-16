package commands

import (
	"context"
	"fmt"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

var DataLengthCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "데이터학습량",
		Description: "봇이 학습한 데ㅇ이터량을 보여줘요.",
	},
	DetailedDescription: DetailedDescription{
		Usage: "/데이터학습량",
	},
	Category: General,
	Flags:    CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	Run: func(ctx *ChatInputContext) error {
		ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		})

		return dataLengthRun(ctx.Inter.Session, ctx.Inter, ctx.Inter.Member.User.Username, ctx.Inter.Member.User.ID)
	},
}

func dataLengthRun(s *discordgo.Session, m any, username, userID string) error {
	textLength, err := databases.GetDatabase().Texts.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return err
	}
	muffinLength, err := databases.GetDatabase().Texts.CountDocuments(context.TODO(), databases.Text{Persona: "muffin"})
	if err != nil {
		return err
	}
	learnLength, err := databases.GetDatabase().Learns.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return err
	}
	userLearnLength, err := databases.GetDatabase().Learns.CountDocuments(context.TODO(), databases.Learn{UserID: userID})
	if err != nil {
		return err
	}

	sum := textLength + learnLength

	return utils.NewMessageSender(m).
		AddComponents(discordgo.Container{
			Components: []discordgo.MessageComponent{
				discordgo.Section{
					Accessory: discordgo.Thumbnail{
						Media: discordgo.UnfurledMediaItem{
							URL: s.State.User.AvatarURL("512"),
						},
					},
					Components: []discordgo.MessageComponent{
						discordgo.TextDisplay{
							Content: fmt.Sprintf("### 저장된 데이터량\n총합: `%d`개", sum),
						},
						discordgo.TextDisplay{
							Content: fmt.Sprintf("- **총 채팅 데이터량**\n> `%d`개", textLength),
						},
						discordgo.TextDisplay{
							Content: fmt.Sprintf("- **머핀 데이터량**\n> `%d`개", muffinLength),
						},
					},
				},
				discordgo.TextDisplay{
					Content: fmt.Sprintf("- **총 지식 데이터량**\n> `%d`개", learnLength),
				},
				discordgo.TextDisplay{
					Content: fmt.Sprintf("- **%s님이 가르쳐준 데이터량**\n> `%d`개", username, userLearnLength),
				},
			},
		}).
		SetComponentsV2(true).
		SetReply(true).
		Send()
}
