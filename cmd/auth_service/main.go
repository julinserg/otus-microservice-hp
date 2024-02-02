package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	auth_app "github.com/julinserg/otus-microservice-hp/internal/auth/app"
	auth_internalhttp "github.com/julinserg/otus-microservice-hp/internal/auth/server/http"
	"github.com/julinserg/otus-microservice-hp/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/auth_config.toml", "Path to configuration file")
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
	}

	f, err := os.OpenFile("auth_service_logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalln("error opening file: " + err.Error())
	}
	defer f.Close()

	logg := logger.New(config.Logger.Level, f)

	logg.Info("auth_service is running...")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	srvAuth := auth_app.New(logg, config.YDisk.ClientId, config.YDisk.ClientSecret)

	endpointHttp := net.JoinHostPort(config.HTTP.Host, config.HTTP.Port)
	serverHttp := auth_internalhttp.NewServer(logg, endpointHttp, srvAuth)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverHttp.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serverHttp.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			return
		}
	}()
	wg.Wait()
}
