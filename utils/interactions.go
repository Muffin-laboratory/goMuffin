package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type ModalData struct {
	CustomId   string                       `json:"custom_id"`
	Title      string                       `json:"title"`
	Components []discordgo.MessageComponent `json:"components"`
}

// InteractionCreate custom data of discordgo.InteractionCreate
type InteractionCreate struct {
	*discordgo.InteractionCreate
	Session *discordgo.Session
	// NOTE: It's only can ApplicationCommand
	Options map[string]*discordgo.ApplicationCommandInteractionDataOption
}

// Reply to this interaction.
func (i *InteractionCreate) Reply(data *discordgo.InteractionResponseData) {
	i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

// GetInteractionOptions to this interaction.
// NOTE: It's only can ApplicationCommand
func GetInteractionOptions(i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optsMap := map[string]*discordgo.ApplicationCommandInteractionDataOption{}
	for _, opt := range i.ApplicationCommandData().Options {
		optsMap[opt.Name] = opt
	}
	return optsMap
}

// DeferReply to this interaction.
func (i *InteractionCreate) DeferReply(ephemeral bool) {
	var flags discordgo.MessageFlags
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}

	i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: flags,
		},
	})
}

// DeferUpdate to this interaction.
func (i *InteractionCreate) DeferUpdate() {
	i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

// EditReply to this interaction.
func (i *InteractionCreate) EditReply(data *discordgo.WebhookEdit) {
	i.Session.InteractionResponseEdit(i.Interaction, data)
}

// Update to this interaction.
func (i *InteractionCreate) Update(data *discordgo.InteractionResponseData) {
	i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: data,
	})
}

func (i *InteractionCreate) ShowModal(data *ModalData) error {
	var reqData struct {
		Type discordgo.InteractionResponseType `json:"type"`
		Data ModalData                         `json:"data"`
	}

	reqData.Type = discordgo.InteractionResponseModal
	reqData.Data = *data
	bin, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(bin)

	req, err := http.NewRequest("POST", discordgo.EndpointInteractionResponse(i.ID, i.Token), buf)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", i.Session.Identify.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := i.Session.Client.Do(req)
	if err != nil {
		return err
	}

	respBin, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", string(respBin))
	}

	defer resp.Body.Close()
	return nil
}
