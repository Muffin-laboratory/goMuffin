package handler

import (
	"log/slog"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
)

func OnApplicationCommandInteractionCreate(i *events.ApplicationCommandInteractionCreate) {
	if err := loader.GetDiscommand().ChatInputRun(i.Data.CommandName(), i); err != nil {
		slog.Error("error in responding chat input command.", "user_id", i.User().ID, "error", err)
		i.CreateMessage(
			discord.NewMessageCreateV2(getErrContainer(i.Client().Rest)).
				WithEphemeral(true),
		)
	}
}

func OnComponentInteractionCreate(i *events.ComponentInteractionCreate) {
	if err := loader.GetDiscommand().ComponentRun(i); err != nil {
		slog.Error("error in responding component.", "user_id", i.User().ID, "error", err)
		i.CreateMessage(
			discord.NewMessageCreateV2(getErrContainer(i.Client().Rest)).
				WithEphemeral(true),
		)
	}
}

func OnAutocompleteInteractionCreate(i *events.AutocompleteInteractionCreate) {
	if err := loader.GetDiscommand().ChatInputAutocomplete(i.Data.CommandName, i); err != nil {
		slog.Error("error in responding autocomplete.", "user_id", i.User().ID, "error", err)
		i.Respond(discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateV2(getErrContainer(i.Client().Rest)).
				WithEphemeral(true),
		)
	}
}

func getErrContainer(client rest.Rest) discord.ContainerComponent {
	owner, _ := client.GetUser(configs.GetConfig().Bot.OwnerID)
	return builders.MakeErrorContainer("오류가 발생하였어요. 만약 계속 발생한다면, `%s`으로 연락해주세요.", owner.Username)
}
