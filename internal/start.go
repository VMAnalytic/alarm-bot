package app

import (
	"gopkg.in/tucnak/telebot.v2"
)

type StartBot struct {
	errCh chan error
}

func NewStartBot(errCh chan error) *StartBot {
	return &StartBot{errCh: errCh}
}

func (s StartBot) Register(tgBot *telebot.Bot) {
	tgBot.Handle(CommandStart, func(m *telebot.Message) {

		var (
			// Inline buttons.
			//
			// Pressing it will cause the client to
			// send the bot a callback.
			//
			// Make sure Unique stays unique as per button kind,
			// as it has to be for callback routing to work.
			//
			selector = &telebot.ReplyMarkup{}

			btnPrev = selector.Data("⬅", "prev", "1111")
			btnNext = selector.Data("➡", "next", "2222")
		)
		menu.Reply(
			menu.Row(btnHelp),
			menu.Row(btnAddContact),
			menu.Row(btnAlarm),
		)
		selector.Inline(
			selector.Row(btnPrev, btnNext),
		)

		_, err := tgBot.Send(m.Sender, "Hello World!", menu)
		if err != nil {
			s.errCh <- err
			//log.Fatal(err)
		}
	})
}

func (s StartBot) Help() string {
	panic("implement me")
}
