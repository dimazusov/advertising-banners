package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimazusov/hw-test/advertising-banners/internal/app"
	"github.com/dimazusov/hw-test/advertising-banners/internal/config"
	"github.com/dimazusov/hw-test/advertising-banners/internal/kafka"
	"github.com/dimazusov/hw-test/advertising-banners/internal/logger"
	internalhttp "github.com/dimazusov/hw-test/advertising-banners/internal/server/http"
	"github.com/dimazusov/hw-test/advertising-banners/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", `confisdgs/config.yaml`, "Path to configuration file")
	flag.Parse()
}

func main() {
	cfg, err := config.New(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	lg, err := logger.New(cfg.Logger.Path, cfg.Logger.Level)
	if err != nil {
		log.Fatalln(err)
	}

	rep, err := storage.NewRepository(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	calendar := app.New(lg, rep, cfg)
	srv := internalhttp.NewServer(cfg, calendar)

	go kafkaListener(cfg)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

		<-signals
		signal.Stop(signals)

		err := srv.Stop(context.Background())
		if err != nil {
			log.Println(err)
		}
	}()

	err = srv.Start(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}

func kafkaListener(cfg *config.Config) {
	consumer, err := kafka.NewConsumer("event", cfg)
	if err != nil {
		log.Println(err)
		return
	}

	err = consumer.Listen(func(b []byte) error {
		fmt.Println(string(b))
		return nil
	})
	if err != nil {
		log.Println(err)
		return
	}
}