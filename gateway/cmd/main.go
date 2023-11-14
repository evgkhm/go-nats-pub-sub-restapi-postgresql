package main

import (
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/config"
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/controller/http"
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/controller/mq/nats"
	"go-nats-pub-sub-restapi-postgresql/gateway/pkg/logging"
	"golang.org/x/net/context"
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

	logger.Info("Config", "http", http.HTTP, "nats", nats.Config)

	nc, js, natsErr := nats.New(nats.Config, logger)
	if natsErr != nil {
		logger.Error("main - nats.New", "err", natsErr)
		return
	}
	logger.Info("nats created and connected")

	natsSubscriber := nats.NewSubscriber(nc, js, logger)
	natsSubscriber.Run(ctx)

	handlers := http.New(js, natsSubscriber, logger)
	r, handlersErr := handlers.InitRoutes()
	if handlersErr != nil {
		logger.Error("main - handlers.InitRoutes", "err", handlersErr)
		return
	}

	runGinErr := r.Run(":" + http.HTTP.Port)
	if runGinErr != nil {
		logger.Error("main - r.Run", "err", runGinErr)
		return
	}
}
