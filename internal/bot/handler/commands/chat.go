package commands

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	subcommands "github.com/Muffin-laboratory/goMuffin/internal/bot/handler/commands/subcommands/chat"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"
)

func init() {
	var chatNameMaxLength = 25

	const (
		name              = "대화"
		chatCommandName   = "하기"
		listCommandName   = "목록"
		createCommandName = "생성"
		deleteCommandName = "삭제"
		setCommandName    = "설정"
	)

	loader.GetDiscommand().RegisterCommand(discord.SlashCommandCreate{
		Name:        name,
		Description: "이 봇이랑 대화해요.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        createCommandName,
				Description: "새로운 채팅을 생성해요.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "이름",
						Description: "채팅방의 이름 (25자 이내)",
						MaxLength:   &chatNameMaxLength,
						Required:    false,
					},
					discord.ApplicationCommandOptionBool{
						Name:        "프롬프트_지정",
						Description: "사용자가 원하는 프롬프트 (모달)",
						Required:    false,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        listCommandName,
				Description: "채팅 목록을 나열해요.",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        chatCommandName,
				Description: "이 봇이랑 대화해요.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "내용",
						Description: "대화할 내용",
						Required:    true,
					},
					discord.ApplicationCommandOptionAttachment{
						Name:        "첨부파일",
						Description: "같이 보낼 파일",
						Required:    false,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        deleteCommandName,
				Description: "채팅을 삭제해요.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:         "이름",
						Description:  "지울 채팅방의 이름",
						MaxLength:    &chatNameMaxLength,
						Required:     true,
						Autocomplete: true,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        setCommandName,
				Description: "봇의 설정을 개인화 해요.",
			},
		},
	})

	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(
			middlewares.CheckIsUserAndBlocked(),
			middlewares.Timeout(loader.Timeout()),
		)

		r.Autocomplete("/"+name, func(e *handler.AutocompleteEvent) error {
			var choices []discord.AutocompleteChoice

			focusedValue := e.Data.Focused().String()

			filter := query.ChatQueryBuilder().SetNameByRegex(focusedValue)
			data, err := repository.GetDatabase().Chats.Find(e.Ctx, filter)
			if err != nil {
				return err
			}

			if len(data) > 25 {
				data = data[:25]
			}

			for _, data := range data {
				choices = append(choices, discord.AutocompleteChoiceString{
					Name:  data.Name,
					Value: data.Name,
				})
			}

			return e.AutocompleteResult(choices)
		})

		r.Route("/"+name, func(r handler.Router) {
			r.Group(func(r handler.Router) {
				r.Use(checkAIMode())
				r.SlashCommand("/"+createCommandName, subcommands.Create)
			})

			r.Group(func(r handler.Router) {
				r.Use(middleware.Defer(discord.InteractionTypeApplicationCommand, false, true))

				r.Group(func(r handler.Router) {
					r.Use(checkAIMode())
					r.Command("/"+listCommandName, subcommands.List)
					r.SlashCommand("/"+deleteCommandName, subcommands.Delete)
				})

				r.SlashCommand("/"+chatCommandName, subcommands.Chat)
				r.Command("/"+setCommandName, subcommands.Settings)
			})
		})
	})
}

func checkAIMode() handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			dbUser, err := repository.GetDatabase().Users.FindByID(e.Ctx, int64(e.User().ID))
			if err != nil {
				return err
			}

			if dbUser.ChattingMode == repository.ChattingMuffinMode {
				container := builders.MakeErrorContainer("채팅모드가 %s여야해요.", repository.ModeString(repository.ChattingAIMode))

				if err := e.CreateMessage(
					discord.NewMessageCreateV2(container).
						WithEphemeral(true),
				); err != nil {
					_, err = e.UpdateInteractionResponse(
						discord.NewMessageUpdateV2([]discord.LayoutComponent{container}),
					)
				}

				return err
			}

			return next(e)
		}
	}
}
