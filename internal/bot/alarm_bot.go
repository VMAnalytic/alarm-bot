package bot

import (
	"context"
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

// AlarmBot Handles alarm button press
//
type AlarmBot struct {
	commonBot
	//client exchange.Client
}

func NewAlarmBot() *AlarmBot {
	return &AlarmBot{}
}

func (b *AlarmBot) Register(ctx context.Context, tgBot *telebot.Bot, errCh chan<- error) {
	b.init(tgBot, errCh)
	tgBot.Handle(&btnAlarm, b.handler(tgBot))
	tgBot.Handle(CommandAlarm, b.handler(tgBot))
}

func (b *AlarmBot) handler(tgBot *telebot.Bot) func(message *telebot.Message) {
	return func(m *telebot.Message) {
		r := &telebot.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}

		r.Reply(r.Row(r.Location("SOS")))
		_, err := tgBot.Send(m.Sender, "Send SOS signal to your friends", r)
		if err != nil {
			b.handleError(err, m.Sender.ID)
			return
		}
	}
}

func (b *AlarmBot) Help() string {
	return fmt.Sprintf("%s - list of all commands", CommandAlarm)
}
