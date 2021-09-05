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
		menu.Inline(
			menu.Row(btnAddContact, btnRemoveContact, btnMyContacts),
			menu.Row(btnMyInfo, btnHelp),
			menu.Row(btnUnsubscribe),
		)

		_, err := tgBot.Send(m.Chat, "Menu", menu)
		if err != nil {
			b.handleError(err, m.Sender.ID)
		}
	})
}

func (b *MenuBot) Help() string {
	panic("implement me")
}
