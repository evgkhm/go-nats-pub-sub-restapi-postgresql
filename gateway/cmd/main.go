package main

import (
	"fmt"
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/config"
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/controller/http"
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/controller/mq/nats"
	"go-nats-pub-sub-restapi-postgresql/gateway/pkg/logging"
	"golang.org/x/net/context"
	"log"
	"log/slog"
	"time"
)

func init() {
	config.InitAll([]config.Config{
		http.HTTPConfig{},
		nats.Conf{},
	})
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger := logging.Logger()
	slog.SetDefault(logger)

	nc, js, err := nats.New(nats.Config, logger)
	if err != nil {
		return
	}
	logger.Info("nats created and connected")

	natsSubscriber := nats.NewSubscriber(nc, js, logger)
	natsSubscriber.Run(ctx)

	handlers := http.New(js, natsSubscriber, logger)
	r := handlers.InitRoutes()

	err = r.Run(":" + http.HTTP.Port)
	if err != nil {
		log.Fatal(fmt.Errorf("main - r.Run: %w", err))
	}
}
