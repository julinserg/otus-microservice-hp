package telegram_bot

import (
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type Storage interface {
	UpdateOrderStatus(idOrder string, status string) error
}

type TelegramBot struct {
	logger        Logger
	token         string
	timeoutUpdate int
}

func New(logger Logger, token string, timeoutUpdate int) *TelegramBot {
	return &TelegramBot{
		logger:        logger,
		token:         token,
		timeoutUpdate: timeoutUpdate,
	}
}
const URI_AUTH_STR = "https://oauth.yandex.ru/authorize?response_type=code&client_id=7c73474791134a019232de6285ca9d34"

func (a *TelegramBot) Start() error {

	a.logger.Info("start bot...")

	bot, err := tgbotapi.NewBotAPI(a.token)
	if err != nil {
		return err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = a.timeoutUpdate

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			switch update.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, URI_AUTH_STR)
				_, err := bot.Send(msg)
				if err != nil {
					a.logger.Error("Bot Send Error: " + err.Error())
				}
			default:
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words only.")
			_, err := bot.Send(msg)
			if err != nil {
				a.logger.Error("Bot Send Error: " + err.Error())
			}
		}
	}
	return nil
}
