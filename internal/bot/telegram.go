package bot

import (
	"context"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
)

type TelegramListener struct {
	tgBot *telebot.Bot
	bots  []Bot
	errCh chan error
}

func NewTelegramListener(tgBot *telebot.Bot, bots []Bot) *TelegramListener {
	return &TelegramListener{tgBot: tgBot, bots: bots, errCh: make(chan error)}
}

func (t *TelegramListener) Listen(ctx context.Context) error {
	errCh := make(chan error, 1)

	for _, b := range t.bots {
		b.Register(ctx, t.tgBot, errCh)
	}

	go t.tgBot.Start()

	for {
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			t.tgBot.Stop()
			return ctx.Err()
		case e := <-errCh:
			if e != nil {
				logrus.Error(e)
			}
		}
	}
}

func (t *TelegramListener) Stop() {
	t.tgBot.Stop()
}
