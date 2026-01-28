package builders

import (
	"github.com/bwmarrin/discordgo"
)

type CommandInteractionOptionsMap map[string]*discordgo.ApplicationCommandInteractionDataOption

type ModalData struct {
	CustomId   string                       `json:"custom_id"`
	Title      string                       `json:"title"`
	Components []discordgo.MessageComponent `json:"components"`
}

type InteractionEdit struct {
	Content         *string                           `json:"content,omitempty"`
	Components      *[]discordgo.MessageComponent     `json:"components,omitempty"`
	Embeds          *[]*discordgo.MessageEmbed        `json:"embeds,omitempty"`
	Flags           *discordgo.MessageFlags           `json:"flags,omitempty"`
	Attachments     *[]*discordgo.MessageAttachment   `json:"attachments,omitempty"`
	AllowedMentions *discordgo.MessageAllowedMentions `json:"allowed_mentions,omitempty"`
}

// InteractionCreate custom data of discordgo.InteractionCreate
type InteractionCreate struct {
	*discordgo.InteractionCreate
	Session *discordgo.Session
	// NOTE: It's only can ApplicationCommand
	Options  CommandInteractionOptionsMap
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

// MakeCommandInteractionOptionsMap to this interaction.
// NOTE: It's only can ApplicationCommand
func MakeCommandInteractionOptionsMap(opts []*discordgo.ApplicationCommandInteractionDataOption) CommandInteractionOptionsMap {
	optsMap := CommandInteractionOptionsMap{}

	for _, opt := range opts {
		optsMap[opt.Name] = opt
	}

	return optsMap
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
func (i *InteractionCreate) EditReply(data *InteractionEdit) error {
	endpoint := discordgo.EndpointWebhookMessage(i.AppID, i.Token, "@original")

	_, err := i.Session.RequestWithBucketID("PATCH", endpoint, *data, discordgo.EndpointWebhookToken("", ""))

	i.Replied = true

	return err
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
func (i *InteractionCreate) ShowModal(data *ModalData) error {
	var reqData struct {
		Type discordgo.InteractionResponseType `json:"type"`
		Data ModalData                         `json:"data"`
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
