package cloud_storage_amqp

import (
	"context"
	"encoding/json"
	"fmt"

	amqp_settings "github.com/julinserg/otus-microservice-hp/internal/amqp/settings"
	amqp_sub "github.com/julinserg/otus-microservice-hp/internal/amqp/sub"
	cloud_storage_app "github.com/julinserg/otus-microservice-hp/internal/cloud_storage/app"
	"github.com/streadway/amqp"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type SrvCloudStorage interface {
	DownloadAndSaveToStorage(cloud_storage_app.FileEvent) error
}

type SrvCloudStorageAMQP struct {
	logger Logger
	srvCS  SrvCloudStorage
	uri    string
}

func New(logger Logger, uri string, srvCS SrvCloudStorage) *SrvCloudStorageAMQP {
	return &SrvCloudStorageAMQP{
		logger: logger,
		uri:    uri,
		srvCS:  srvCS,
	}
}

func (a *SrvCloudStorageAMQP) StartReceive(ctx context.Context) error {
	conn, err := amqp.Dial(a.uri)
	if err != nil {
		return err
	}
	c := amqp_sub.New("SrvCloudStorageAMQP", conn, a.logger)
	msgs, err := c.Consume(ctx, amqp_settings.QueueFileTransfer, amqp_settings.ExchangeFileTransfer,
		"direct", "")
	if err != nil {
		return err
	}

	a.logger.Info("start consuming file event...")

	for m := range msgs {
		fileEvent := cloud_storage_app.FileEvent{}
		json.Unmarshal(m.Data, &fileEvent)
		if err != nil {
			return err
		}
		a.logger.Info(fmt.Sprintf("receive new message:%+v\n", fileEvent))

		err := a.srvCS.DownloadAndSaveToStorage(fileEvent)
		if err != nil {
			a.logger.Warn("Error SrvCloudStorageAMQP: " + err.Error())
		}
	}
	return nil
}
