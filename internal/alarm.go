package app

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

type AlarmBot struct {
	//client exchange.Client
}

func NewAlarmBot() *AlarmBot {
	return &AlarmBot{}
}

func (b *AlarmBot) Register(tgBot *telebot.Bot) {
	tgBot.Handle("/alarm", b.handler(tgBot))
}

func (b *AlarmBot) handler(tgBot *telebot.Bot) func(message *telebot.Message) {
	return func(m *telebot.Message) {
		r := &telebot.ReplyMarkup{ResizeReplyKeyboard: true}
		r.Location("Send location")
		r.URL("Visit", "https://google.com")
		r.Text("V11111sits")
		tgBot.Send(m.Sender, "sdsd", r)
	}
}

//func (s StartBot) handler() func()  {
//	return func(tgBot *telebot.Bot) {
//		tgBot.Send(m.Sender, "Hello World!", menu)
//	}
//}
//
//func (s StartBot) handler2() func()  {
//	return func(tgBot *telebot.Bot) {
//		tgBot.Send(m.Sender, "Hello World!", menu)
//	}
//}

func (b AlarmBot) Help() string {
	return fmt.Sprintf("%s - list of all commands", CommandAlarm)
}
