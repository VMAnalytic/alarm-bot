package bot

import (
	"context"
	"fmt"
	app "github.com/VMAnalytic/alarm-bot/internal"
	"github.com/VMAnalytic/alarm-bot/pkg/convertor"
	"github.com/pkg/errors"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

type TextBot struct {
	sessionStorage app.SessionStorage
	userStorage    app.UserStorage
	commonBot
}

const (
	separator = ","
)

func NewTextBot(sessionStorage app.SessionStorage, userStorage app.UserStorage) *TextBot {
	return &TextBot{sessionStorage: sessionStorage, userStorage: userStorage}
}

func (b *TextBot) Register(ctx context.Context, tgBot *telebot.Bot, errCh chan<- error) {
	b.init(tgBot, errCh)

	//b.tgBot.Handle()
	b.tgBot.Handle(telebot.OnText, func(m *telebot.Message) {
		var uid = m.Sender.ID

		//Handles add contacts case
		if b.sessionStorage.ExistInState(ctx, uid, app.AddingContacts) {
			user, err := b.userStorage.Get(ctx, uid)
			if err != nil {
				b.handleError(err, uid)
				return
			}
			ids := strings.Split(m.Text, separator)
			for _, id := range ids {
				id, err := convertor.ToInt(id)
				if err != nil {
					b.handleError(errors.New("unknown id"), uid)
					return
				}

				e, err := b.userStorage.Exists(ctx, id)
				if err != nil {
					b.handleError(err, uid)
					return
				}

				if !e {
					b.handleError(app.NewErrContactNotFound(id), uid)
					return
				}

				err = user.AddContact(app.NewContact(id))
				if err != nil {
					b.handleError(err, uid)
					return
				}
			}
			err = b.userStorage.Add(ctx, user)
			if err != nil {
				b.handleError(err, uid)
				return
			}

			b.sessionStorage.Delete(ctx, uid)

			_, err = tgBot.Send(m.Sender, "Contacts added")
			if err != nil {
				b.handleError(err, uid)
				return
			}

			return
		}

		//Handles remove contacts case
		if b.sessionStorage.ExistInState(ctx, uid, app.RemovingContacts) {
			user, err := b.userStorage.Get(ctx, uid)
			if err != nil {
				b.handleError(err, uid)
				return
			}
			ids := strings.Split(m.Text, separator)

			//Add all contacts
			for _, id := range ids {
				id, err := convertor.ToInt(id)
				if err != nil {
					b.handleError(errors.New("unknown id"), uid)
					return
				}

				e, err := b.userStorage.Exists(ctx, id)
				if err != nil {
					b.handleError(err, uid)
					return
				}

				if !e {
					b.handleError(errors.Wrap(app.ErrNotFound, fmt.Sprintf("Contact with ID: %v", id)), uid)
					return
				}

				err = user.RemoveContact(app.NewContact(id))
				if err != nil {
					b.handleError(err, uid)
					return
				}
			}
			err = b.userStorage.Add(ctx, user)
			if err != nil {
				b.handleError(err, uid)
			}

			b.sessionStorage.Delete(ctx, uid)

			_, err = tgBot.Send(m.Sender, "Contacts removed")
			if err != nil {
				b.handleError(err, uid)
				return
			}

			return
		}

		_, err := tgBot.Send(m.Sender, "Unexpected message!", mainMenu)
		if err != nil {
			b.errCh <- errors.New("unexpected message")
		}
	})
}

func (b *TextBot) Help() string {
	panic("implement me")
}
