package main

import (
	"context"
	"github.com/sirupsen/logrus/hooks/writer"
	"log"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"github.com/VMAnalytic/alarm-bot/config"
)

func init()  {
	//logrus.SetOutput(io.Discard)
	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
	})
	logrus.AddHook(&writer.Hook{
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})
}

func main() {
	ctx := context.TODO()
	cfg, err := alar
	if err != nil {
		logrus.Fatal(err)
	}

	tgBot, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		//URL: "http://195.129.111.17:8012",

		Token:  cfg.Telegram.Token,
		Poller: &tb.LongPoller{Timeout: time.Duration(cfg.Telegram.Timeout) * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	tgListener := NewTelegramLis

	err = tgListener.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
