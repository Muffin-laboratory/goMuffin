package commands

import (
	"context"
	"fmt"
	"log"
	"strings"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

type chatCommandType string

var (
	chatCommandChatting chatCommandType = "하기"
	chatCommandList     chatCommandType = "목록"
	chatCommandCreate   chatCommandType = "생성"
	chatCommandDelete   chatCommandType = "삭제"
)

var ChatCommand *Command = &Command{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "대화",
		Description: "이 봇이랑 대화해요.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandList),
				Description: "채팅 목록을 나열해요.",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandCreate),
				Description: "새로운 채팅을 생성해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "이름",
						Description: "채팅의 이름을 정해요. (25자 이내)",
						MaxLength:   25,
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandChatting),
				Description: "이 봇이랑 대화해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "내용",
						Description: "대화할 내용",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        string(chatCommandDelete),
				Description: "채팅을 삭제해요.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "제목",
						Description: "채팅의 제목",
						Required:    true,
					},
				},
			},
		},
	},
	Aliases: []string{"채팅"},
	DetailedDescription: &DetailedDescription{
		Usage: fmt.Sprintf("%s대화 (목록/생성) [채팅 이름]", configs.Config.Bot.Prefix),
		Examples: []string{
			fmt.Sprintf("%s대화 목록", configs.Config.Bot.Prefix),
			fmt.Sprintf("%s대화 생성 머핀 냠냠", configs.Config.Bot.Prefix),
		},
	},
	Category:                   Chatting,
	RegisterApplicationCommand: true,
	RegisterMessageCommand:     true,
	Flags:                      CommandFlagsIsRegistered | CommandFlagsIsBlocked,
	ChatInputRun: func(ctx *ChatInputContext) error {
		ctx.Inter.DeferReply(nil)

		var cType chatCommandType
		var str string

		if opt, ok := ctx.Inter.Options[string(chatCommandChatting)]; ok {
			ctx.Inter.DeferReply(nil)
			cType = chatCommandChatting
			str = opt.Options[0].StringValue()
		} else if opt, ok := ctx.Inter.Options[string(chatCommandCreate)]; ok {
			ctx.Inter.DeferReply(&discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			})
			cType = chatCommandCreate
			str = opt.Options[0].StringValue()
		} else if _, ok := ctx.Inter.Options[string(chatCommandList)]; ok {
			ctx.Inter.DeferReply(nil)
			cType = chatCommandList
		} else if _, ok := ctx.Inter.Options[string(chatCommandDelete)]; ok {
			ctx.Inter.DeferReply(nil)
			cType = chatCommandDelete
			str = opt.Options[0].StringValue()
		}
		return chatCommandRun(cType, ctx.Inter, ctx.Inter.User, str)
	},
	MessageRun: func(ctx *MsgContext) error {
		if len((*ctx.Args)) < 1 {
			goto RequiredValue
		}

		switch (*ctx.Args)[0] {
		case string(chatCommandCreate):
			if len((*ctx.Args)) < 2 {
				return utils.NewMessageSender(ctx.Msg).
					AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "채팅방의 이름을 정해야해요."})).
					SetComponentsV2(true).
					SetReply(true).
					Send()
			}

			name := strings.Trim(strings.Join((*ctx.Args), " "), " ")
			if len([]rune(name)) > 25 {
				return utils.NewMessageSender(ctx.Msg).
					AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "채팅방의 이름은 25자를 초과할 수 없어요."})).
					SetComponentsV2(true).
					SetReply(true).
					Send()
			}

			return chatCommandRun(chatCommandCreate, ctx.Msg, ctx.Msg.Author, name)
		case string(chatCommandList):
			return chatCommandRun(chatCommandList, ctx.Msg, ctx.Msg.Author, "")
		case string(chatCommandDelete):
			if len((*ctx.Args)) < 2 {
				return utils.NewMessageSender(ctx.Msg).
					AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "채팅방의 이름을 적어야해요."})).
					SetComponentsV2(true).
					SetReply(true).
					Send()
			}
			return chatCommandRun(chatCommandDelete, ctx.Msg, ctx.Msg.Author, strings.Join((*ctx.Args)[1:], " "))
		default:
			goto RequiredValue
		}

	RequiredValue:
		return utils.NewMessageSender(ctx.Msg).
			AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "명령어의 첫번째 인자는 `생성`, `목록`중에 하나여야 해요."})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	},
}

func chatCommandRun(cType chatCommandType, m any, user *discordgo.User, contentOrName string) error {
	switch cType {
	// 채팅하기는 슬래시 커맨드만 가능
	case chatCommandChatting:
		i := m.(*utils.InteractionCreate)

		str, err := chatbot.ChatBot.GetResponse(user, contentOrName)
		if err != nil {
			log.Println(err)
			i.EditReply(&utils.InteractionEdit{
				Content: &str,
			})
			return nil
		}

		result := chatbot.ParseResult(str, i.Session, i)
		return i.EditReply(&utils.InteractionEdit{
			Content: &result,
		})
	case chatCommandCreate:
		_, err := databases.CreateChat(user.ID, contentOrName)
		if err != nil {
			return err
		}
		return utils.NewMessageSender(m).
			AddComponents(utils.GetSuccessContainer(discordgo.TextDisplay{Content: fmt.Sprintf("%s를 생성했어요. 이제 현재 채팅은 %s에요.", contentOrName, contentOrName)})).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	case chatCommandList:
		var dbUser databases.User
		var data []databases.Chat
		var sections []discordgo.Section
		var containers []*discordgo.Container

		cur, err := databases.Database.Chats.Find(context.TODO(), databases.Chat{UserId: user.ID})
		if err != nil {
			return err
		}

		err = cur.All(context.TODO(), &data)
		if err != nil {
			return err
		}

		err = databases.Database.Users.FindOne(context.TODO(), databases.User{UserId: user.ID}).Decode(&dbUser)
		if err != nil {
			return err
		}

		if len(data) == 0 {
			return utils.NewMessageSender(m).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "채팅이 단 하나도 없어요. 새로운 채팅을 만들거나, 대화를 시작해 채팅을 만들어주세요."})).
				SetComponentsV2(true).
				SetReply(true).
				SetEphemeral(true).
				Send()
		}

		for i, data := range data {
			var isDisabled bool
			var textDisplay discordgo.TextDisplay

			if data.Id == dbUser.ChatId {
				textDisplay = discordgo.TextDisplay{
					Content: fmt.Sprintf("**%d. %s\n (선택됨)**", i+1, data.Name),
				}

				isDisabled = true
			} else {
				textDisplay = discordgo.TextDisplay{
					Content: fmt.Sprintf("%d. %s\n", i+1, data.Name),
				}

				isDisabled = false
			}

			sections = append(sections, discordgo.Section{
				Accessory: discordgo.Button{
					Label:    "선택",
					Style:    discordgo.SuccessButton,
					CustomID: utils.MakeSelectChat(data.Id.Hex(), i+1, user.ID),
					Disabled: isDisabled,
				},
				Components: []discordgo.MessageComponent{textDisplay},
			})
		}

		textDisplay := discordgo.TextDisplay{Content: fmt.Sprintf("### %s님의 채팅목록", user.GlobalName)}
		container := &discordgo.Container{Components: []discordgo.MessageComponent{textDisplay}}
		for i, section := range sections {
			container.Components = append(container.Components, section, discordgo.Separator{})

			if (i+1)%10 == 0 {
				containers = append(containers, container)
				container = &discordgo.Container{Components: []discordgo.MessageComponent{textDisplay}}
				continue
			}
		}

		if len(container.Components) > 1 {
			containers = append(containers, container)
		}

		return utils.PaginationEmbedBuilder(m).
			AddContainers(containers...).
			Start()
	case chatCommandDelete:
		var data []databases.Chat

		cur, err := databases.Database.Chats.Find(context.TODO(), databases.Chat{Name: contentOrName})
		if err != nil {
			return err
		}

		err = cur.All(context.TODO(), &data)
		if err != nil {
			return err
		}

		if len(data) == 0 {
			return utils.NewMessageSender(m).
				AddComponents(utils.GetErrorContainer(discordgo.TextDisplay{Content: "해당하는 채팅을 찾을 수 없어요."})).
				SetComponentsV2(true).
				SetReply(true).
				Send()
		}

		if len(data) > 1 {
			var sections []discordgo.Section
			var containers []*discordgo.Container

			for i, data := range data {
				sections = append(sections, discordgo.Section{
					Accessory: discordgo.Button{
						Label:    "삭제",
						Style:    discordgo.DangerButton,
						CustomID: utils.MakeDeleteChat(data.Id.Hex(), i+1, user.ID),
					},
					Components: []discordgo.MessageComponent{
						discordgo.TextDisplay{
							Content: fmt.Sprintf("%d. %s\n", i+1, data.Name),
						},
					},
				})
			}

			textDisplay := discordgo.TextDisplay{Content: fmt.Sprintf("### %s님의 채팅목록\n- **주의: 이 채팅방을 삭제하면 이 채팅방의 내역을 다시는 못 써요.**", user.GlobalName)}
			container := &discordgo.Container{Components: []discordgo.MessageComponent{textDisplay}}
			for i, section := range sections {
				container.Components = append(container.Components, section, discordgo.Separator{})

				if (i+1)%10 == 0 {
					containers = append(containers, container)
					container = &discordgo.Container{Components: []discordgo.MessageComponent{textDisplay}}
					continue
				}
			}

			if len(container.Components) > 1 {
				containers = append(containers, container)
			}

			return utils.PaginationEmbedBuilder(m).
				AddContainers(containers...).
				Start()
		}

		return utils.NewMessageSender(m).
			AddComponents(discordgo.Container{
				Components: []discordgo.MessageComponent{
					discordgo.TextDisplay{Content: fmt.Sprintf("### 채팅 %s 삭제", contentOrName)},
					discordgo.TextDisplay{Content: "- **주의: 이 채팅방을 삭제하면 이 채팅방의 내역을 다시는 못 써요.**"},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "삭제",
								Style:    discordgo.DangerButton,
								CustomID: utils.MakeDeleteChat(data[0].Id.Hex(), 0, user.ID),
							},
							discordgo.Button{
								Label:    "취소",
								Style:    discordgo.PrimaryButton,
								CustomID: utils.MakeDeleteChatCancel(user.ID),
							},
						},
					},
				},
			}).
			SetComponentsV2(true).
			SetReply(true).
			Send()
	}
	return nil
}
