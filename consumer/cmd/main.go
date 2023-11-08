package main

import (
	"fmt"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/config"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/controller/mq/nats"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/repository/postgres"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/transactions"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/usecase"
	"go-nats-pub-sub-restapi-postgresql/gateway/pkg/logging"
	"log"
)

func init() {
	config.InitAll([]config.Config{
		postgres.Conf{},
		nats.Conf{},
	})
}

func main() {
	logger := logging.Logger()

	logger.Info("Config", "postgres", postgres.Config, "nats", nats.Config)

	postgresDB, err := postgres.NewPostgresDB(postgres.Config, logger)
	if err != nil {
		return
	}

	postgresRepository := postgres.New(postgresDB)

	txService := transactions.New(postgresDB)

	useCases := usecase.New(postgresRepository, txService, logger)

	nc, js, err := nats.New(nats.Config, logger)
	if err != nil {
		log.Fatal(fmt.Errorf("main - nats.New: %w", err))
	}

	natsSubscriber := nats.NewSubscriber(nc, js, useCases, logger)
	natsSubscriber.Run()
}
