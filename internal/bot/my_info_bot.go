package bot

import (
	"context"
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

type MyInfoBot struct {
	errCh chan error
}

func NewMyInfoBot() *MyInfoBot {
	return &MyInfoBot{errCh: make(chan error)}
}

func (b *MyInfoBot) Register(ctx context.Context, tgBot *telebot.Bot, errCh chan<- error) {
	tgBot.Handle(CommandGetMyInfo, b.handler(tgBot))
	//tgBot.Handle(&btnMyInfo, b.addContactsHandler(tgBot))
	tgBot.Handle(&btnMyInfo, func(c *telebot.Callback) {
		b.handler(tgBot)(c.Message)

		//res := telebot.CallbackResponse{Text: "Insert id"}
		//tgBot.Send(c.Sender, "hello")
		//tgBot.Respond(c, &res)
		//b.Respond(c, &tb.CallbackResponse{...})
	})
}

func (b *MyInfoBot) handler(tgBot *telebot.Bot) func(message *telebot.Message) {
	return func(m *telebot.Message) {
		id := m.Sender.ID
		fName := m.Sender.FirstName
		lName := m.Sender.LastName
		userName := m.Sender.Username

		//TODO: format msg
		msg := fmt.Sprintf(
			"ID: %-20v \nFirst name: %-20s \nLast Name: %-20s \nUserName: %-20s",
			id,
			fName,
			lName,
			userName,
		)
		_, err := tgBot.Send(m.Sender, msg)
		if err != nil {
			b.errCh <- err
		}
	}
}

func (b *MyInfoBot) Help() string {
	panic("implement me")
}
