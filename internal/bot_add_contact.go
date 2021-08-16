package app

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

type AddContactBot struct {
	//client exchange.Client
}

func NewAddContactBot() *AddContactBot {
	return &AddContactBot{}
}

func (b *AddContactBot) Register(tgBot *telebot.Bot) {
	tgBot.Handle("/add_contact", b.handler(tgBot))
}

func (b *AddContactBot) handler(tgBot *telebot.Bot) func(message *telebot.Message) {
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

func (b *AddContactBot) Help() string {
	return fmt.Sprintf("%s - list of all commands", CommandAddContact)
}
