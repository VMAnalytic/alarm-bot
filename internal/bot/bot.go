package bot

import (
	"context"
	"fmt"
	app "github.com/VMAnalytic/alarm-bot/internal"
	"github.com/pkg/errors"
	"gopkg.in/tucnak/telebot.v2"
)

//List of available commands
const (
	CommandStart         = "/start"
	CommandMenu          = "/menu"
	CommandAddContact    = "/add_contacts"
	CommandRemoveContact = "/remove_contacts"
	CommandMyContacts    = "/my_contacts"
	CommandGetMyInfo     = "/my_info"
	CommandAlarm         = "/alarm"
	CommandUnsubscribe   = "/unsubscribe"
)

var (
	mainMenu         = &telebot.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}
	menu             = &telebot.ReplyMarkup{}
	btnAlarm         = mainMenu.Location(fmt.Sprintf("%s ALARM!", []byte{0xF0, 0x9F, 0x86, 0x98}))
	btnAddContact    = mainMenu.Data(fmt.Sprintf("%s Add contacts", []byte{0xF0, 0x9F, 0x93, 0x92}), "add_contacts")
	btnRemoveContact = mainMenu.Data(fmt.Sprintf("%s Remove contacts", []byte{0xF0, 0x9F, 0x93, 0x92}), "remove_contacts")
	btnMyContacts    = mainMenu.Data(fmt.Sprintf("%s My contact list", []byte{0xF0, 0x9F, 0x93, 0x92}), CommandMyContacts)
	btnMyInfo        = mainMenu.Data(fmt.Sprintf("%s Get my info", []byte{0xF0, 0x9F, 0x93, 0x83}), "info")
	btnUnsubscribe   = menu.Data(fmt.Sprintf("%s Delete me!", []byte{0xF0, 0x9F, 0x93, 0x83}), CommandUnsubscribe)
	btnHelp          = mainMenu.Text("ℹ Help")
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
		b.errCh <- errors.Wrap(err, fmt.Sprintf("error for user: %v", toUserID))
		var (
			errContNotFound *app.ErrContactNotFound
		)

		switch true {
		case errors.As(err, &errContNotFound):
			_, _ = b.tgBot.Send(&telebot.User{ID: toUserID}, err.Error())
		case errors.Is(err, app.ErrNotFound) || errors.Is(err, app.ErrContactsLimitExceeded):
			_, _ = b.tgBot.Send(&telebot.User{ID: toUserID}, err.Error())
		default:
			_, _ = b.tgBot.Send(&telebot.User{ID: toUserID}, "Unexpected error occurred.")
		}
	}
}
