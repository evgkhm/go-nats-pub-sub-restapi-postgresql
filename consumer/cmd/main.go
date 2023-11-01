package main

import (
	"fmt"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/config"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/controller/mq/nats"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/repository/postgres"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/transactions"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/usecase"
	"golang.org/x/net/context"
	"log"
	"time"
)

func init() {
	config.InitAll([]config.Config{
		postgres.Conf{},
		nats.Conf{},
	})
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	postgresDB, err := postgres.NewPostgresDB(postgres.Config)
	if err != nil {
		log.Fatal(fmt.Errorf("main - postgres.NewPostgresDB: %w", err))
	}

	postgresRepository := postgres.New(postgresDB)

	txService := transactions.New(postgresDB)

	useCases := usecase.New(postgresRepository, txService)

	nc, js, err := nats.New(nats.Config)
	if err != nil {
		log.Fatal(fmt.Errorf("main - nats.New: %w", err))
	}

	natsSubscriber := nats.NewSubscriber(nc, js, useCases)
	natsSubscriber.Run(ctx)

	//handlers := http2.New(useCases, js)
	//r := handlers.InitRoutes()
	//
	//err = r.Run(":" + http2.HTTP.Port)
	//if err != nil {
	//	log.Fatal(fmt.Errorf("main - r.Run: %w", err))
	//}
}
