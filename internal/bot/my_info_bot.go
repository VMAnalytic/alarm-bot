package bot

import (
	"context"
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

type MyInfoBot struct {
	commonBot
}

func NewMyInfoBot() *MyInfoBot {
	return &MyInfoBot{}
}

func (b *MyInfoBot) Register(ctx context.Context, tgBot *telebot.Bot, errCh chan<- error) {
	b.init(tgBot, errCh)

	tgBot.Handle(CommandGetMyInfo, b.handler)
	tgBot.Handle(&btnMyInfo, func(c *telebot.Callback) {
		c.Message.Sender = c.Sender //change sender from bot to original user
		b.handler()(c.Message)
		err := b.tgBot.Respond(c, &telebot.CallbackResponse{}) //empty response
		if err != nil {
			b.handleError(err, c.Sender.ID)
			return
		}
	})
}

func (b *MyInfoBot) handler() func(message *telebot.Message) {
	return func(m *telebot.Message) {
		id := m.Sender.ID
		fName := m.Sender.FirstName
		lName := m.Sender.LastName
		userName := m.Sender.Username

		//TODO: format msg
		msg := fmt.Sprintf(
			"*ID*: %-20v \n*First name*: %-20s \n*Last Name*: %-20s \n*UserName*: %-20s",
			id,
			fName,
			lName,
			userName,
		)
		_, err := b.tgBot.Send(m.Chat, msg, telebot.ModeMarkdown)
		if err != nil {
			b.errCh <- err
		}
	}
}

func (b *MyInfoBot) Help() string {
	panic("implement me")
}
