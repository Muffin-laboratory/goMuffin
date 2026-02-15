package components

import (
	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader/middlewares"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func init() {
	loader.GetDiscommand().RegisterHandler(func(r handler.Router) {
		r.Use(
			middlewares.CheckIDMiddleware(),
			middlewares.TimeoutAndDeferMiddleware(loader.Timeout(), discord.InteractionTypeComponent, false, true),
		)

		r.Component(utils.SelectChat+"/{command}", func(e *handler.ComponentEvent) error {
			var sections []discord.SectionComponent
			var containers []discord.ContainerComponent

			command := e.Vars["command"]

			filter := query.KnowledgeQueryBuilder().SetUserID(int64(e.User().ID)).SetCommand(command)
			data, err := repository.GetDatabase().Knowledge.Find(e.Ctx, filter)
			if err != nil {
				return err
			}

			if len(data) == 0 {
				_, err := e.UpdateInteractionResponse(
					discord.NewMessageUpdateV2([]discord.LayoutComponent{
						builders.MakeErrorContainer("해당 결과를 찾을 수 없어요."),
					}),
				)
				return err
			}

			for _, data := range data {
				sections = append(sections,
					discord.NewSection(
						discord.NewTextDisplayf("**%s**\n", data.Result),
					).
						WithAccessory(
							discord.NewDangerButton("삭제", utils.MakeDeleteKnowledge(data.ID.Hex(), e.User().ID.String())),
						),
				)
			}

			textDisplay := discord.NewTextDisplayf("### %s에 대한 목록", command)
			container := discord.NewContainer().WithComponents(textDisplay)
			for i, section := range sections {
				container = container.AddComponents(section, discord.NewSeparator(discord.SeparatorSpacingSizeSmall))

				if (i+1)%10 == 0 {
					containers = append(containers, container)
					container = container.WithComponents(textDisplay)
					continue
				}
			}

			if len(container.Components) > 1 {
				containers = append(containers, container)
			}

			return builders.NewPaginatedContainer(e.User().ID, true).
				AddContainers(containers...).
				Start(e)
		})
	})
}
