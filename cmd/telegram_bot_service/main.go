package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/julinserg/otus-microservice-hp/internal/logger"
	telegram_bot_amqp "github.com/julinserg/otus-microservice-hp/internal/telegram_bot/amqp"
	telegram_bot_app "github.com/julinserg/otus-microservice-hp/internal/telegram_bot/app"
	telegram_bot "github.com/julinserg/otus-microservice-hp/internal/telegram_bot/bot"
	telegram_bot_imitation_internalhttp "github.com/julinserg/otus-microservice-hp/internal/telegram_bot/server/http"
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
		value, _ = os.LookupEnv("USC_AUTH_URI")
		config.AuthSrv.URI = value
		value, _ = os.LookupEnv("USC_DEBUG_HOST")
		config.Debug.Host = value
		value, _ = os.LookupEnv("USC_DEBUG_PORT")
		config.Debug.Port = value
	}

	f, err := os.OpenFile("telegram_bot_service_logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalln("error opening file: " + err.Error())
	}
	defer f.Close()

	logg := logger.New(config.Logger.Level, f)

	botMQ := telegram_bot_amqp.New(logg, config.AMQP.URI)

	srvBot := telegram_bot_app.New(logg, config.AuthSrv.URI, botMQ)

	tb := telegram_bot.New(logg, config.TGBot.Token, config.TGBot.Timeout, srvBot)

	endpointHttp := net.JoinHostPort(config.Debug.Host, config.Debug.Port)
	serverHttp := telegram_bot_imitation_internalhttp.NewServer(logg, endpointHttp, srvBot)

	logg.Info("telegram_bot_service is running...")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := tb.Start(); err != nil {
			logg.Error("failed to start bot: " + err.Error())
			return
		}
	}()
	go func() {
		defer wg.Done()
		if err := serverHttp.Start(); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			return
		}
	}()
	wg.Wait()
}
