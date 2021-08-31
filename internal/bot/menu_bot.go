package bot

import (
	"context"
	"gopkg.in/tucnak/telebot.v2"
)

type MenuBot struct {
	commonBot
}

func NewMenuBot() *MenuBot {
	return &MenuBot{}
}

func (b *MenuBot) Register(ctx context.Context, tgBot *telebot.Bot, errCh chan<- error) {
	b.init(tgBot, errCh)
	tgBot.Handle(CommandMenu, func(m *telebot.Message) {

		menu.Reply(
			menu.Row(btnAddContact),
			menu.Row(btnRemoveContact),
			menu.Row(btnMyInfo, btnHelp),
		)

		_, err := tgBot.Send(m.Sender, "Welcome to alarm bot!", menu)
		if err != nil {
			b.handleError(err, m.Sender.ID)
		}
	})
}

func (b *MenuBot) Help() string {
	panic("implement me")
}
