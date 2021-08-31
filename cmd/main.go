package main

import (
	"context"
	firebase "firebase.google.com/go"
	app "github.com/VMAnalytic/alarm-bot/internal"
	"github.com/VMAnalytic/alarm-bot/internal/bot"
	"github.com/VMAnalytic/alarm-bot/internal/storage"
	"google.golang.org/api/option"
	"log"
	"time"

	"github.com/VMAnalytic/alarm-bot/config"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
	})
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	ctx := context.TODO()
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	tgBot, err := tb.NewBot(tb.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &tb.LongPoller{Timeout: time.Duration(cfg.Telegram.Timeout) * time.Second},
	})

	if err != nil {
		logrus.Fatal(err)
	}

	var (
		sessionStorage app.SessionStorage
		userStorage    app.UserStorage
	)

	if cfg.Env.IsGoogle() {
		conf := &firebase.Config{ProjectID: cfg.Google.ProjectID}
		fApp, err := firebase.NewApp(ctx, conf, option.WithCredentialsFile("secret.json"))
		if err != nil {
			logrus.Fatal(err)
		}

		firestoreClient, err := fApp.Firestore(ctx)

		if err != nil {
			logrus.Fatal(err)
		}

		userStorage = storage.NewUserFirestoreStorage(firestoreClient)
	}

	sessionStorage = storage.NewMemorySessionStorage()

	tgListener := bot.NewTelegramListener(tgBot, []bot.Bot{
		bot.NewStartBot(userStorage),
		bot.NewAlarmBot(),
		bot.NewAddContactBot(sessionStorage, userStorage),
		bot.NewMyInfoBot(),
		bot.NewTextBot(sessionStorage, userStorage),
		bot.NewLocationBot(sessionStorage, userStorage),
	})

	if err := tgListener.Listen(ctx); err != nil {
		log.Fatalf("Telegram listener failed, %v", err)
	}
}
