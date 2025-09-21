package knowledge

import (
	"context"
	"fmt"
	"slices"

	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func List(i *utils.InteractionCreate, opts utils.CommandInteractionOptionsMap) error {
	var command string
	var data []databases.Learn
	var items []string
	var sections []*discordgo.Section
	var containers []*discordgo.Container

	if opt, ok := opts["단어"]; ok {
		command = opt.StringValue()
	}

	filter := bson.D{{Key: "user_id", Value: i.User.ID}}

	if command != "" {
		filter = append(filter, bson.E{Key: "command", Value: bson.M{"$regex": command}})
	} else {
		command = "전체"
	}

	cur, err := databases.GetDatabase().Learns.Find(context.TODO(), filter)
	if err != nil {
		return err
	}

	defer cur.Close(context.TODO())

	if err = cur.All(context.TODO(), &data); err != nil {
		return err
	}

	if len(data) == 0 {
		return utils.NewMessageSender(i).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "해당 결과를 찾을 수 없어요."})).
			SetComponentsV2(true).
			SetEphemeral(true).
			Send()
	}

	for _, data := range data {
		if slices.Contains(items, data.Command) {
			continue
		}

		items = append(items, data.Command)
	}

	for _, item := range items {
		sections = append(sections, &discordgo.Section{
			Accessory: discordgo.Button{
				CustomID: utils.MakeSelectKnowledge(item),
				Label:    "자세히 보기",
				Style:    discordgo.PrimaryButton,
			},
			Components: []discordgo.MessageComponent{
				discordgo.TextDisplay{
					Content: fmt.Sprintf("- **%s**", item),
				},
			},
		})
	}

	title := &discordgo.Section{
		Accessory: discordgo.Thumbnail{Media: discordgo.UnfurledMediaItem{URL: i.User.AvatarURL("512")}},
		Components: []discordgo.MessageComponent{
			discordgo.TextDisplay{Content: fmt.Sprintf("### %s님에 대한 지식", i.User.GlobalName)},
			discordgo.TextDisplay{Content: fmt.Sprintf("- 총 `%d`개", len(items))},
			discordgo.TextDisplay{Content: fmt.Sprintf("- %s에 대한 검색 결과", command)},
		},
	}
	container := &discordgo.Container{Components: []discordgo.MessageComponent{title, discordgo.Separator{}}}

	for i, section := range sections {
		container.Components = append(container.Components, section, discordgo.Separator{})

		if (i+1)%5 == 0 {
			containers = append(containers, container)
			container = &discordgo.Container{Components: []discordgo.MessageComponent{title, discordgo.Separator{}}}
			continue
		}
	}

	if len(container.Components) > 1 {
		containers = append(containers, container)
	}

	return utils.PaginationContainerBuilder(i).
		AddContainers(containers...).
		Start()
}
