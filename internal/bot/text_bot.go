package bot

import (
	"context"
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

func NewTextBot(sessionStorage app.SessionStorage, userStorage app.UserStorage) *TextBot {
	return &TextBot{sessionStorage: sessionStorage, userStorage: userStorage}
}

func (b *TextBot) Register(ctx context.Context, tgBot *telebot.Bot, errCh chan<- error) {
	b.init(tgBot, errCh)

	b.tgBot.Handle(telebot.OnText, func(m *telebot.Message) {
		var uid = m.Sender.ID

		//Handles add contacts case
		if b.sessionStorage.ExistInState(ctx, uid, app.AddingContacts) {
			user, err := b.userStorage.Get(ctx, uid)
			if err != nil {
				b.handleError(err, uid)
			}
			ids := strings.Split(m.Text, ",")
			for _, id := range ids {
				id, err := convertor.ToInt(id)
				if err != nil {
					b.handleError(errors.New("unknown id"), uid)
				}

				e, err := b.userStorage.Exists(ctx, id)
				if err != nil {
					b.handleError(err, uid)
				}

				if !e {
					b.handleError(app.ErrNotFound, uid)
				}

				err = user.AddContact(app.NewContact(id))
				if err != nil {
					b.handleError(err, uid)
				}
			}
			err = b.userStorage.Add(ctx, user)
			if err != nil {
				b.handleError(err, uid)
			}

			b.sessionStorage.Delete(ctx, uid)

			return
		}

		//Handles remove contacts case
		if b.sessionStorage.ExistInState(ctx, uid, app.RemovingContacts) {
			user, err := b.userStorage.Get(ctx, uid)
			if err != nil {
				b.handleError(err, uid)
			}
			ids := strings.Split(m.Text, ",")
			for _, id := range ids {
				id, err := convertor.ToInt(id)
				if err != nil {
					b.handleError(errors.New("unknown id"), uid)
				}

				e, err := b.userStorage.Exists(ctx, id)
				if err != nil {
					b.handleError(err, uid)
				}

				if !e {
					b.handleError(app.ErrNotFound, uid)
				}

				err = user.RemoveContact(app.NewContact(id))
				if err != nil {
					b.handleError(err, uid)
				}
			}
			err = b.userStorage.Add(ctx, user)
			if err != nil {
				b.handleError(err, uid)
			}

			b.sessionStorage.Delete(ctx, uid)

			return
		}

		_, err := tgBot.Send(m.Sender, "Unexpected message!", menu)
		if err != nil {
			b.errCh <- errors.New("unexpected message")
		}
	})
}

func (b *TextBot) Help() string {
	panic("implement me")
}
