package builders

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

type MessageSender struct {
	Embeds          []discord.Embed
	Content         string
	Components      []discord.LayoutComponent
	Ephemeral       bool
	Reply           bool
	ComponentsV2    bool
	AllowedMentions *discord.AllowedMentions
	m               any
}

type CommandCreate struct {
	*events.ApplicationCommandInteractionCreate
	Responded bool
}

func NewMessageSender(m any) *MessageSender {
	return &MessageSender{m: m}
}

func (s *MessageSender) AddEmbeds(embeds ...discord.Embed) *MessageSender {
	s.Embeds = append(s.Embeds, embeds...)
	return s
}

func (s *MessageSender) AddComponents(components ...discord.LayoutComponent) *MessageSender {
	s.Components = append(s.Components, components...)
	return s
}

func (s *MessageSender) SetContent(content string) *MessageSender {
	s.Content = content
	return s
}

func (s *MessageSender) SetEphemeral(ephemeral bool) *MessageSender {
	s.Ephemeral = ephemeral
	return s
}

func (s *MessageSender) SetReply(reply bool) *MessageSender {
	s.Reply = reply
	return s
}

func (s *MessageSender) SetAllowedMentions(allowedMentions discord.AllowedMentions) *MessageSender {
	s.AllowedMentions = &allowedMentions
	return s
}

func (s *MessageSender) SetComponentsV2(componentsV2 bool) *MessageSender {
	s.ComponentsV2 = componentsV2
	return s
}

func (s *MessageSender) Send() error {
	switch m := s.m.(type) {
	case *events.MessageCreate:
		_, err := m.Client().Rest.CreateMessage(
			m.ChannelID,
			discord.NewMessageCreate().
				WithContent(s.Content).
				AddEmbeds(s.Embeds...).
				AddComponents(s.Components...).
				WithIsComponentsV2(s.ComponentsV2).
				WithAllowedMentions(s.AllowedMentions).
				WithMessageReference(m.Message.MessageReference),
		)

		return err
	case *CommandCreate:
		if m.Responded {
			_, err := m.Client().Rest.UpdateInteractionResponse(
				m.ApplicationID(),
				m.Token(),
				discord.NewMessageUpdate().
					WithContent(s.Content).
					AddEmbeds(s.Embeds...).
					AddComponents(s.Components...).
					WithIsComponentsV2(s.ComponentsV2),
			)
			return err
		}

		err := m.CreateMessage(
			discord.NewMessageCreate().
				WithContent(s.Content).
				AddEmbeds(s.Embeds...).
				AddComponents(s.Components...).
				WithIsComponentsV2(s.ComponentsV2).
				WithEphemeral(s.Ephemeral),
		)
		if err != nil {
			return err
		}

		m.Responded = true
	case *handler.CommandEvent:
		return m.CreateMessage(
			discord.NewMessageCreate().
				WithContent(s.Content).
				AddEmbeds(s.Embeds...).
				AddComponents(s.Components...).
				WithIsComponentsV2(s.ComponentsV2).
				WithEphemeral(s.Ephemeral),
		)
	case *handler.InteractionEvent:
		return m.CreateMessage(
			discord.NewMessageCreate().
				WithContent(s.Content).
				AddEmbeds(s.Embeds...).
				AddComponents(s.Components...).
				WithIsComponentsV2(s.ComponentsV2).
				WithEphemeral(s.Ephemeral),
		)
	}
	return nil
}
