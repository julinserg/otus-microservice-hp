package telegram_bot

import (
	"fmt"
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type SrvBot interface {
	GetAuthRequestString() (string, error)
	SendFileEvent(url string, chatId int64) error
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type TelegramBot struct {
	logger        Logger
	tokenBot      string
	timeoutUpdate int
	srvBot        SrvBot
}

func New(logger Logger, tokenBot string, timeoutUpdate int, srvBot SrvBot) *TelegramBot {
	return &TelegramBot{
		logger:        logger,
		tokenBot:      tokenBot,
		timeoutUpdate: timeoutUpdate,
		srvBot:        srvBot,
	}
}

func (a *TelegramBot) Start() error {

	a.logger.Info("start bot...")

	bot, err := tgbotapi.NewBotAPI(a.tokenBot)
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
				authStr, err := a.srvBot.GetAuthRequestString()
				if err != nil {
					a.logger.Error("Get Auth String Error: " + err.Error())
				}
				authStr = authStr + "&state=%d"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(authStr, update.Message.Chat.ID))
				_, err = bot.Send(msg)
				if err != nil {
					a.logger.Error("Bot Send Error: " + err.Error())
				}
			default:
			}
		} else if update.Message.Photo != nil {
			fileDescript := (*update.Message.Photo)[len(*update.Message.Photo)-1]
			url, err := bot.GetFileDirectURL(fileDescript.FileID)
			if err != nil {
				a.logger.Error("GetFileDirectURL Error: " + err.Error())
			}
			err = a.srvBot.SendFileEvent(url, update.Message.Chat.ID)
			if err != nil {
				a.logger.Error("SendFileEvent Error: " + err.Error())
			}
		}
	}
	return nil
}
