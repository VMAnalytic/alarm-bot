package bot

import (
	"context"
	"fmt"
	app "github.com/VMAnalytic/alarm-bot/internal"
	"github.com/VMAnalytic/alarm-bot/pkg/convertor"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

type ContactsBot struct {
	userStorage    app.UserStorage
	sessionStorage app.SessionStorage
	commonBot
}

func NewAddContactBot(sessionStorage app.SessionStorage, userStorage app.UserStorage) *ContactsBot {
	return &ContactsBot{sessionStorage: sessionStorage, userStorage: userStorage}
}

func (b *ContactsBot) Register(ctx context.Context, tgBot *telebot.Bot, errCh chan<- error) {
	b.init(tgBot, errCh)

	tgBot.Handle(CommandAddContact, b.addContactsHandler(ctx))
	tgBot.Handle(&btnAddContact, b.addContactsHandler(ctx))

	tgBot.Handle(CommandRemoveContact, b.removeContactsHandler(ctx))
	tgBot.Handle(&btnAddContact, b.removeContactsHandler(ctx))

	tgBot.Handle(CommandMyContacts, b.myContactsHandler(ctx))
	tgBot.Handle(&btnAddContact, b.myContactsHandler(ctx))
}

func (b *ContactsBot) addContactsHandler(ctx context.Context) func(message *telebot.Message) {
	return func(m *telebot.Message) {
		var (
			msg             string
			uid             = m.Sender.ID
			currentContacts string
			err             error
		)

		user, err := b.userStorage.Get(ctx, uid)
		if err != nil {
			b.handleError(err, uid)
			return
		}

		s := app.NewSession(uid)
		s.Transition(app.AddingContacts)
		b.sessionStorage.Add(ctx, s)

		currentContacts = b.getContactList(user)

		msg = "Please add up to *5* contacts to your alarm list. \n" +
			"Enter the list of your friends IDs. Use comma to separate IDs. \n" +
			"*Example*: 1234,5678,9012,3456 \n" +
			fmt.Sprintf("Your current contacts list is: %s", currentContacts)

		_, err = b.tgBot.Send(m.Chat, msg, telebot.ModeMarkdown)

		if err != nil {
			b.handleError(err, uid)
		}
	}
}

func (b *ContactsBot) removeContactsHandler(ctx context.Context) func(message *telebot.Message) {
	return func(m *telebot.Message) {
		var (
			msg             string
			uid             = m.Sender.ID
			currentContacts string
			err             error
		)

		user, err := b.userStorage.Get(ctx, uid)
		if err != nil {
			b.handleError(err, uid)
			return
		}

		s := app.NewSession(uid)
		s.Transition(app.RemovingContacts)
		b.sessionStorage.Add(ctx, s)

		currentContacts = b.getContactList(user)

		msg = fmt.Sprintf("Remove contacts from your alarm list. \n"+
			"Enter the list of your friends IDs. Use comma to separate IDs. \n"+
			"Your current contacts list is: %s", currentContacts)

		_, err = b.tgBot.Send(m.Chat, msg)

		if err != nil {
			b.handleError(err, uid)
		}
	}
}

func (b *ContactsBot) myContactsHandler(ctx context.Context) func(message *telebot.Message) {
	return func(m *telebot.Message) {
		var (
			msg             string
			uid             = m.Sender.ID
			currentContacts string
			err             error
		)

		user, err := b.userStorage.Get(ctx, uid)
		if err != nil {
			b.handleError(err, uid)
			return
		}

		currentContacts = b.getContactList(user)

		msg = fmt.Sprintf("Your current contacts list is: %s", currentContacts)

		_, err = b.tgBot.Send(m.Chat, msg)

		if err != nil {
			b.handleError(err, uid)
		}
	}
}

func (b *ContactsBot) getContactList(u *app.User) string {
	if len(u.Contacts) == 0 {
		return "EMPTY"
	}

	contacts := strings.Builder{}
	for _, contact := range u.Contacts {
		contacts.WriteString("\n- " + convertor.ToString(contact.UserID))
	}

	return contacts.String()
}

func (b *ContactsBot) Help() string {
	return fmt.Sprintf("%s - list of all commands", CommandAddContact)
}
