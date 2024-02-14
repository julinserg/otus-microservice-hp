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

	cloud_storage_amqp "github.com/julinserg/otus-microservice-hp/internal/cloud_storage/amqp"
	cloud_storage_app "github.com/julinserg/otus-microservice-hp/internal/cloud_storage/app"
	cloud_storage_debug_internalhttp "github.com/julinserg/otus-microservice-hp/internal/cloud_storage/server/http"
	"github.com/julinserg/otus-microservice-hp/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/cloud_storage_config.toml", "Path to configuration file")
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
		value, _ = os.LookupEnv("USC_AUTH_URI")
		config.AuthSrv.URI = value
		value, _ = os.LookupEnv("USC_STORAGE_FOLDER")
		config.Storage.Folder = value
		value, _ = os.LookupEnv("USC_DEBUG_TOKEN")
		config.Debug.TokenYD = value
	}

	f, err := os.OpenFile("cloud_storage_service_logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalln("error opening file: " + err.Error())
	}
	defer f.Close()

	logg := logger.New(config.Logger.Level, f)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	srvCS := cloud_storage_app.New(logg, config.AuthSrv.URI, ctx, config.Debug.TokenYD, config.Storage.Folder)

	csMQ := cloud_storage_amqp.New(logg, config.AMQP.URI, srvCS)

	endpointHttp := net.JoinHostPort(config.Debug.Host, config.Debug.Port)
	serverHttp := cloud_storage_debug_internalhttp.NewServer(logg, endpointHttp, srvCS)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverHttp.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("cloud_storage_service is running...")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := csMQ.StartReceive(ctx); err != nil {
			logg.Error("failed to start MQ worker(order): " + err.Error())
			cancel()
			return
		}
	}()
	go func() {
		defer wg.Done()
		if err := serverHttp.Start(); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			return
		}
	}()
	wg.Wait()
}
