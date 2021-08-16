package app

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

type command string

//Commands
const (
	CommandStart = "/start"
	CommandAddContact     = "/add_contact"
	CommandAlarm    = "/alarm"
)

//Formatting
const (
	Bold   = "*"
	Italic = "_"
	Mono   = "`"
)

var (
	menu        = &telebot.ReplyMarkup{OneTimeKeyboard: true}
	btnAlarm    = menu.Location(fmt.Sprintf("%s ALARM!", []byte{0xE2, 0x9D, 0x97}))
	btnAddContact     = menu.Text("ℹ Add contact")
	btnHelp     = menu.Text("ℹ CommandHelp")
)

type handlerFunc func(message *telebot.Message)

type Bot interface {
	Register(tgBot *telebot.Bot)
	Help() string
}
