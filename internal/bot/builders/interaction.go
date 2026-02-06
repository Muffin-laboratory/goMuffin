package builders

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// InteractionCreate custom data of discordgo.InteractionCreate
type InteractionCreate struct {
	*discordgo.InteractionCreate
	Session *discordgo.Session
	Ctx     context.Context
	// NOTE: It's only can ApplicationCommand
	Deferred bool
	Replied  bool
}

// Reply to this interaction.
func (i *InteractionCreate) Reply(data *discordgo.InteractionResponseData) error {
	if err := i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	}); err != nil {
		return err
	}

	i.Replied = true
	return nil
}

// NOTE: It's only can ApplicationCommand
func GetInteractionUser(i *discordgo.InteractionCreate) *discordgo.User {
	if i.Member != nil {
		return i.Member.User
	}

	if i.User != nil {
		return i.User
	}

	return nil
}

// DeferReply to this interaction.
func (i *InteractionCreate) DeferReply(data *discordgo.InteractionResponseData) error {
	if err := i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: data,
	}); err != nil {
		return err
	}

	i.Deferred = true

	return nil
}

// FetchReply gets message that was sent.
func (i *InteractionCreate) FetchReply() (*discordgo.Message, error) {
	return i.Session.WebhookMessage(i.AppID, i.Token, "@original")
}

// DeferUpdate to this interaction.
func (i *InteractionCreate) DeferUpdate() error {
	if err := i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	}); err != nil {
		return err
	}

	i.Deferred = true

	return nil
}

// EditReply to this interaction.
func (i *InteractionCreate) EditReply(data *discordgo.WebhookEdit) error {
	if _, err := i.Session.WebhookMessageEdit(i.AppID, i.Token, "@original", data); err != nil {
		return err
	}

	i.Replied = true

	return nil
}

// Update to this interaction.
func (i *InteractionCreate) Update(data *discordgo.InteractionResponseData) error {
	if err := i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: data,
	}); err != nil {
		return err
	}

	i.Replied = true

	return nil
}

// ShowModal shows modal to this interaction.
func (i *InteractionCreate) ShowModal(data *Modal) error {
	var reqData struct {
		Type discordgo.InteractionResponseType `json:"type"`
		Data Modal                             `json:"data"`
	}

	reqData.Type = discordgo.InteractionResponseModal
	reqData.Data = *data

	endpoint := discordgo.EndpointInteractionResponse(i.ID, i.Token)
	_, err := i.Session.RequestWithBucketID("POST", endpoint, reqData, endpoint)

	return err
}

func (i *InteractionCreate) Autocomplete(options []*discordgo.ApplicationCommandOptionChoice) error {
	if err := i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: options,
		},
	}); err != nil {
		return err
	}

	return nil
}
