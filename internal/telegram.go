package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

type TelegramListener struct {
	tgBot *telebot.Bot
	bots  []Bot
	errCh chan error
}

func NewTelegramListener(tgBot *telebot.Bot, bots []Bot) *TelegramListener {
	return &TelegramListener{tgBot: tgBot, bots: bots, errCh: make(chan error)}
}

func (t TelegramListener) Listen(ctx context.Context) <-chan error {
	var errCh chan error
	defer close(errCh)

	go func() {
		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				t.tgBot.Stop()
			case err := <- errCh:
				logrus.Errorf("%s", err)
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	for _, b := range t.bots {
		b.Register(t.tgBot)
	}

	go t.tgBot.Start()

	return t.errCh

	err := <-errCh

	if err != nil {
		t.tgBot.Stop()
		return err
	}

	return nil
}
