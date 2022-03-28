package entities

import "github.com/andersfylling/disgord"

type ComponentBuilder struct {
	*disgord.MessageComponent
}

func (cb *ComponentBuilder) Select(placeholder, customID string, options []*disgord.SelectMenuOption) *ComponentBuilder {
	selectMenu := &disgord.MessageComponent{
		Type:        disgord.MessageComponentSelectMenu,
		Placeholder: placeholder,
		CustomID:    customID,
		Options:     options,
		MaxValues:   1,
	}
	cb.Components = append(cb.Components, selectMenu)
	return cb
}

func (cb *ComponentBuilder) Button(style disgord.ButtonStyle, label, customID string, opts ...*disgord.MessageComponent) *ComponentBuilder {
	btn := &disgord.MessageComponent{
		Type:     disgord.MessageComponentButton,
		Style:    style,
		Label:    label,
		CustomID: customID,
	}
	if len(opts) > 0 {
		opt := opts[0]
		btn.Url = opt.Url
		btn.Emoji = opt.Emoji
		btn.Disabled = opt.Disabled
	}
	cb.Components = append(cb.Components, btn)
	return cb
}

type MsgBuilder struct {
	*disgord.CreateInteractionResponse
}

func (m *MsgBuilder) Res() *disgord.CreateInteractionResponse {
	return m.CreateInteractionResponse
}

func (m *MsgBuilder) Component(cb *ComponentBuilder) *MsgBuilder {
	m.Data.Components = append(m.Data.Components, cb.MessageComponent)
	return m
}

func (m *MsgBuilder) Embed(embed *disgord.Embed) *MsgBuilder {
	m.Data.Embeds = []*disgord.Embed{embed}
	return m
}

func (m *MsgBuilder) ResType(t disgord.InteractionCallbackType) *MsgBuilder {
	m.Type = t
	return m
}

func (m *MsgBuilder) Content(content string) *MsgBuilder {
	m.Data.Content = content
	return m
}

func CreateComponent() *ComponentBuilder {
	return &ComponentBuilder{
		MessageComponent: &disgord.MessageComponent{
			Type: disgord.MessageComponentActionRow,
		},
	}
}

func CreateMsg() *MsgBuilder {
	return &MsgBuilder{
		CreateInteractionResponse: &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{},
		},
	}
}
