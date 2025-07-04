package core

import (
	"fmt"

	"github.com/lugvitc/steve/ext"
	"github.com/lugvitc/steve/ext/context"
	"github.com/lugvitc/steve/ext/handlers"
	waLogger "github.com/lugvitc/steve/logger"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func whois(client *whatsmeow.Client, ctx *context.Context) error {
	msg := ctx.Message.Message
	if msg.Message.ExtendedTextMessage != nil && msg.Message.ExtendedTextMessage.ContextInfo != nil {
		var jid types.JID
		if msg.Message.ExtendedTextMessage.ContextInfo.Participant != nil {
			var err error
			jid, err = types.ParseJID(*msg.Message.ExtendedTextMessage.ContextInfo.Participant)
			if err != nil {
				fmt.Println(err.Error())
				_, _ = reply(client, ctx.Message, fmt.Sprintf("failed to get user: %s", err.Error()))
				return ext.EndGroups
			}
		} else {
			jid = msg.Info.MessageSource.Chat
		}
		users, err := client.GetUserInfo([]types.JID{jid})
		if err != nil {
			_, _ = reply(client, ctx.Message, fmt.Sprintf("failed to get user: %s", err.Error()))
			return ext.EndGroups
		}
		text := fmt.Sprintf(`*User Info*
*ID*: %s
*Server*: %s
*About*: %s
*Verified*: %s
`, jid.String(), jid.Server, func() string {
			if users[jid].Status != "" {
				return users[jid].Status
			}
			return "None"
		}(), func() string {
			if users[jid].VerifiedName != nil {
				return "```true```"
			}
			return "```false```"
		}())
		_, _ = reply(client, ctx.Message, text)
		return ext.EndGroups
	}
	if ctx.Message.Message.Info.IsGroup {
		return ext.EndGroups
	}
	jid := ctx.Message.Info.Chat
	users, err := client.GetUserInfo([]types.JID{jid})
	if err != nil {
		_, _ = reply(client, ctx.Message, fmt.Sprintf("failed to get user: %s", err.Error()))
		return ext.EndGroups
	}
	text := fmt.Sprintf(`*User Info*
*ID*: %s
*Server*: %s
*About*: %s
*Verified*: %s
`, jid.String(), jid.Server, func() string {
		if users[jid].Status != "" {
			return users[jid].Status
		}
		return "None"
	}(), func() string {
		if users[jid].VerifiedName != nil {
			return "```true```"
		}
		return "```false```"
	}())
	_, _ = reply(client, ctx.Message, text)
	return ext.EndGroups
}

func devices(client *whatsmeow.Client, ctx *context.Context) error {
	msg := ctx.Message.Message
	if msg.Message.ExtendedTextMessage != nil && msg.Message.ExtendedTextMessage.ContextInfo != nil {
		var jid types.JID
		if msg.Message.ExtendedTextMessage.ContextInfo.Participant != nil {
			var err error
			jid, err = types.ParseJID(*msg.Message.ExtendedTextMessage.ContextInfo.Participant)
			if err != nil {
				fmt.Println(err.Error())
				_, _ = reply(client, ctx.Message, fmt.Sprintf("failed to get user: %s", err.Error()))
				return ext.EndGroups
			}
		} else {
			jid = msg.Info.MessageSource.Chat
		}
		devices, err := client.GetUserDevices([]types.JID{jid})
		if err != nil {
			_, _ = reply(client, ctx.Message, fmt.Sprintf("failed to get user: %s", err.Error()))
			return ext.EndGroups
		}
		text := "*Devices*:"
		for _, device := range devices {
			text += fmt.Sprintf("\n- %s", device.String())
		}
		text += fmt.Sprintf("\n\nTotal signed in devices: %d", len(devices))
		_, _ = reply(client, ctx.Message, text)
		return ext.EndGroups
	}
	if ctx.Message.Message.Info.IsGroup {
		return ext.EndGroups
	}
	jid := ctx.Message.Info.Chat
	devices, err := client.GetUserDevices([]types.JID{jid})
	if err != nil {
		_, _ = reply(client, ctx.Message, fmt.Sprintf("failed to get user: %s", err.Error()))
		return ext.EndGroups
	}
	text := "*Devices*:"
	for _, device := range devices {
		text += fmt.Sprintf("\n- %s", device.String())
	}
	text += fmt.Sprintf("\n\nTotal signed in devices: %d", len(devices))
	_, _ = reply(client, ctx.Message, text)
	return ext.EndGroups
}

func (*Module) LoadWhois(dispatcher *ext.Dispatcher) {
	ppLogger := LOGGER.Create("whois")
	defer ppLogger.Println("Loaded Whois module")
	dispatcher.AddHandler(
		handlers.NewCommand("whois", authorizedOnly(whois), ppLogger.Create("whois-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("Get details of a user."),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("devices", authorizedOnly(devices), ppLogger.Create("devices-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription("List devices being used by a user."),
	)
}
