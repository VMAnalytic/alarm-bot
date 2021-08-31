package bot

import (
	"context"
	"fmt"
	app "github.com/VMAnalytic/alarm-bot/internal"
	"gopkg.in/tucnak/telebot.v2"
)

//Commands
const (
	CommandStart         = "/start"
	CommandMenu          = "/menu"
	CommandAddContact    = "/add_contacts"
	CommandRemoveContact = "/remove_contacts"
	CommandMyContacts    = "/my_contacts"
	CommandGetMyInfo     = "/get_info"
	CommandAlarm         = "/alarm"
	CommandUnsubscribe   = "/unsubscribe"
)

var (
	menu             = &telebot.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}
	btnAlarm         = menu.Location(fmt.Sprintf("%s ALARM!", []byte{0xF0, 0x9F, 0x86, 0x98}))
	btnAddContact    = menu.Data(fmt.Sprintf("%s Add contacts", []byte{0xF0, 0x9F, 0x93, 0x92}), CommandAddContact)
	btnRemoveContact = menu.Data(fmt.Sprintf("%s Remove contacts", []byte{0xF0, 0x9F, 0x93, 0x92}), CommandRemoveContact)
	btnMyInfo        = menu.Data(fmt.Sprintf("%s Get my info", []byte{0xF0, 0x9F, 0x93, 0x83}), "info")
	btnHelp          = menu.Text("ℹ Help")
)

type Bot interface {
	Register(ctx context.Context, tgBot *telebot.Bot, errCh chan<- error)
	Help() string
}

type commonBot struct {
	tgBot *telebot.Bot
	errCh chan<- error
}

func (b *commonBot) init(tgBot *telebot.Bot, errCh chan<- error) {
	b.tgBot = tgBot
	b.errCh = errCh
}

func (b *commonBot) handleError(err error, toUserID int) {
	if err != nil {
		b.errCh <- err
		switch err {
		case app.ErrContactAlreadyExists:
		default:
			_, _ = b.tgBot.Send(&telebot.User{ID: toUserID}, "Error")
		}
	}
}
