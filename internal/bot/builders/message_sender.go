package builders

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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
			discord.NewMessageCreateBuilder().
				SetContent(s.Content).
				AddEmbeds(s.Embeds...).
				AddComponents(s.Components...).
				SetIsComponentsV2(s.ComponentsV2).
				SetAllowedMentions(s.AllowedMentions).
				SetMessageReference(m.Message.MessageReference).
				Build(),
		)

		return err
	case *CommandCreate:
		if m.Responded {
			_, err := m.Client().Rest.UpdateInteractionResponse(
				m.ApplicationID(),
				m.Token(),
				discord.NewMessageUpdateBuilder().
					SetContent(s.Content).
					AddEmbeds(s.Embeds...).
					AddComponents(s.Components...).
					SetIsComponentsV2(s.ComponentsV2).
					Build(),
			)
			return err
		}

		err := m.CreateMessage(
			discord.NewMessageCreateBuilder().
				SetContent(s.Content).
				AddEmbeds(s.Embeds...).
				AddComponents(s.Components...).
				SetIsComponentsV2(s.ComponentsV2).
				SetEphemeral(s.Ephemeral).
				Build(),
		)
		if err != nil {
			return err
		}

		m.Responded = true
	}
	return nil
}
