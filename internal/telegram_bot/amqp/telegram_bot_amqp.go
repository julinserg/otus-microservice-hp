package telegram_bot_amqp

import (
	"encoding/json"
	"fmt"

	amqp_pub "github.com/julinserg/otus-microservice-hp/internal/amqp/pub"
	amqp_settings "github.com/julinserg/otus-microservice-hp/internal/amqp/settings"
	telegram_bot_app "github.com/julinserg/otus-microservice-hp/internal/telegram_bot/app"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type SrvBotAMQP struct {
	logger Logger
	pub    amqp_pub.AmqpPub
	uri    string
}

func New(logger Logger, uri string) *SrvBotAMQP {
	return &SrvBotAMQP{
		logger: logger,
		pub:    *amqp_pub.New(logger),
		uri:    uri,
	}
}

func (a *SrvBotAMQP) PublishFileEvent(event telegram_bot_app.FileEvent) error {
	eventStr, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := a.pub.Publish(a.uri, amqp_settings.ExchangeFileTransfer, "direct",
		"", string(eventStr), true); err != nil {
		return err
	}
	a.logger.Info(fmt.Sprintf("publish order for queue is OK ( chatID: %d URL: %s)", event.ChatId, event.URL))
	return nil
}
