package bot

import (
	"context"
	app "github.com/VMAnalytic/alarm-bot/internal"
	"gopkg.in/tucnak/telebot.v2"
)

type UnsubscribeBot struct {
	userStorage app.UserStorage
	commonBot
}

func NewUnsubscribeBot(userStorage app.UserStorage) *UnsubscribeBot {
	return &UnsubscribeBot{userStorage: userStorage}
}

func (b *UnsubscribeBot) Register(ctx context.Context, tgBot *telebot.Bot, errCh chan<- error) {
	b.init(tgBot, errCh)

	// On inline button pressed (callback)
	b.tgBot.Handle(&btnUnsubscribe, func(c *telebot.Callback) {
		msg := "Are u sure?"

		tgBot.Send(c.Sender, "Welcome to alarm bot!", mainMenu)

		// Always respond!
		b.tgBot.Respond(c, &telebot.CallbackResponse{Text: msg, ShowAlert: true})
	})

	tgBot.Handle(CommandUnsubscribe, func(m *telebot.Message) {
		var uid = m.Sender.ID

		err := b.userStorage.Remove(ctx, m.Sender.ID)
		if err != nil {
			b.handleError(err, uid)
			return
		}

		_, err = tgBot.Send(m.Sender, "You were deleted")
		if err != nil {
			b.handleError(err, m.Sender.ID)
		}
	})
}

func (b *UnsubscribeBot) Help() string {
	panic("implement me")
}
