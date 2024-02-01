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

	cloud_storage_internalhttp "github.com/julinserg/otus-microservice-hp/internal/cloud_storage/server/http"
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

	/*yaDisk, err := yadisk.NewYaDisk(ctx, http.DefaultClient, &yadisk.Token{AccessToken: config.YDisk.Token})
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
	}*/

	endpointHttp := net.JoinHostPort(config.HTTP.Host, config.HTTP.Port)
	serverHttp := cloud_storage_internalhttp.NewServer(logg, endpointHttp, config.YDisk.ClientSecret)

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
