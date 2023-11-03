package main

import (
	"fmt"
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/config"
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/controller/http"
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/controller/mq/nats"
	"golang.org/x/net/context"
	"log"
	"time"
)

func init() {
	config.InitAll([]config.Config{
		http.HTTPConfig{},
		nats.Conf{},
	})
}

func main() {
	//TODO: add logger
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nc, js, err := nats.New(nats.Config)
	if err != nil {
		log.Fatal(fmt.Errorf("main - nats.New: %w", err))
	}

	natsSubscriber := nats.NewSubscriber(nc, js)
	natsSubscriber.Run(ctx)

	handlers := http.New(js, natsSubscriber)
	r := handlers.InitRoutes()

	err = r.Run(":" + http.HTTP.Port)
	if err != nil {
		log.Fatal(fmt.Errorf("main - r.Run: %w", err))
	}
}
