package delivery

import (
	"butler/application/commands/helper"
	"butler/constants"
	"strings"

	cartHandler "butler/application/domains/cart/delivery/discord/handler"
	pickHandler "butler/application/domains/pick/delivery/discord/handler"
	makersuiteHandler "butler/application/domains/promt_ai/makersuite/handler"

	"github.com/bwmarrin/discordgo"
)

type Handler interface {
	GetCommandsHandler(*discordgo.Session, *discordgo.MessageCreate)
}

type commandHandler struct {
	discord           *discordgo.Session
	makersuiteHandler makersuiteHandler.Handler
	cartHandler       cartHandler.Handler
	pickHandler       pickHandler.Handler
}

func NewCommandHandler(
	discord *discordgo.Session,
	makersuiteHandler makersuiteHandler.Handler,
	cartHandler cartHandler.Handler,
	pickHandler pickHandler.Handler,
) Handler {
	return &commandHandler{
		discord:           discord,
		makersuiteHandler: makersuiteHandler,
		cartHandler:       cartHandler,
		pickHandler:       pickHandler,
	}
}

func (c *commandHandler) GetCommandsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, constants.BOT_COMMAND_PREFIX) && !helper.CheckMention(m, s.State.User) {
		return
	}

	var err error
	switch {
	case helper.CheckPrefixCommand(m.Content, constants.COMMAND_HELP):
		err = helper.HandleHelpCommand(s, m)
	case helper.CheckMention(m, s.State.User):
		err = c.makersuiteHandler.Ask(s, m)
	case helper.CheckPrefixCommand(m.Content, constants.COMMAND_RESET_CART):
		err = c.cartHandler.ResetCart(s, m)
	case helper.CheckPrefixCommand(m.Content, constants.COMMAND_READY_OUTBOUND):
		err = c.pickHandler.ReadyPickOutbound(s, m)
	}
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
	}
}
