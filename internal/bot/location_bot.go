package bot

import (
	"context"
	app "github.com/VMAnalytic/alarm-bot/internal"
	"golang.org/x/sync/errgroup"
	"gopkg.in/tucnak/telebot.v2"
)

type LocationBot struct {
	userStorage    app.UserStorage
	sessionStorage app.SessionStorage
	commonBot
}

func NewLocationBot(sessionStorage app.SessionStorage, userStorage app.UserStorage) *LocationBot {
	return &LocationBot{userStorage: userStorage, sessionStorage: sessionStorage}
}

func (b *LocationBot) Register(ctx context.Context, tgBot *telebot.Bot, errCh chan<- error) {
	b.init(tgBot, errCh)
	tgBot.Handle(telebot.OnLocation, func(m *telebot.Message) {
		var uid = m.Sender.ID
		user, err := b.userStorage.Get(ctx, uid)
		if err != nil {
			b.handleError(err, uid)
			return
		}

		//TODO: add retries
		g, _ := errgroup.WithContext(ctx)
		for _, contact := range user.Contacts {
			c := contact
			g.Go(func() error {
				_, err := b.tgBot.Send(&telebot.User{ID: c.UserID}, "Help NEEDED")
				if err != nil {
					return err
				}

				_, err = b.tgBot.Send(&telebot.User{ID: c.UserID}, m.Location)
				if err != nil {
					return err
				}

				return nil
			})
		}

		err = g.Wait()
		if err != nil {
			b.handleError(err, uid)
		}
	})
}

func (b *LocationBot) Help() string {
	panic("implement me")
}
