package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/julinserg/otus-microservice-hp/internal/logger"
	telegram_bot_app "github.com/julinserg/otus-microservice-hp/internal/telegram_bot/app"
	telegram_bot "github.com/julinserg/otus-microservice-hp/internal/telegram_bot/bot"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/telegram_bot_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	err := config.Read(configFile)
	if err != nil {
		var value string
		value, _ = os.LookupEnv("USC_LOG_LEVEL")
		config.Logger.Level = value
		value, _ = os.LookupEnv("USC_AMQP_URI")
		config.AMQP.URI = value
		value, _ = os.LookupEnv("USC_TGBOT_TOKEN")
		config.TGBot.Token = value
		value, _ = os.LookupEnv("USC_TGBOT_TIMEOUT")
		timeout, _ := strconv.Atoi(value)
		config.TGBot.Timeout = timeout
	}

	f, err := os.OpenFile("telegram_bot_service_logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalln("error opening file: " + err.Error())
	}
	defer f.Close()

	logg := logger.New(config.Logger.Level, f)

	srvBot := telegram_bot_app.New(logg, config.AuthSrv.URI)

	tb := telegram_bot.New(logg, config.TGBot.Token, config.TGBot.Timeout, srvBot)
	logg.Info("telegram_bot_service is running...")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := tb.Start(); err != nil {
			logg.Error("failed to start bot: " + err.Error())
			return
		}
	}()
	wg.Wait()
}
