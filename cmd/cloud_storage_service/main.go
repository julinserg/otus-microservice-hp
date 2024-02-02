package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julinserg/otus-microservice-hp/internal/logger"
	yadisk "github.com/nikitaksv/yandex-disk-sdk-go"
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
	}

	f, err := os.OpenFile("cloud_storage_service_logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalln("error opening file: " + err.Error())
	}
	defer f.Close()

	logg := logger.New(config.Logger.Level, f)

	logg.Info("cloud_storage_service is running...")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	yaDisk, err := yadisk.NewYaDisk(ctx, http.DefaultClient, &yadisk.Token{AccessToken: config.YDisk.Token})
	if err != nil {
		panic(err.Error())
	}
	disk, err := yaDisk.GetDisk([]string{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("TotalSpace", disk.TotalSpace)
	fmt.Println("UsedSpace", disk.UsedSpace)
	l, err := yaDisk.GetFlatFilesList([]string{}, 10, "", 0, false, "", "")
	if err != nil {
		panic(err.Error())
	}
	for _, item := range l.Items {
		fmt.Println("Name", item.Name)
	}
}
