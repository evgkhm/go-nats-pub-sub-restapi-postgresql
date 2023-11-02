package main

import (
	"fmt"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/config"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/controller/mq/nats"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/repository/postgres"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/transactions"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/usecase"
	"log"
)

func init() {
	config.InitAll([]config.Config{
		postgres.Conf{},
		nats.Conf{},
	})
}

func main() {
	//ctx, cancel := context.WithCancel(context.Background())
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

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
	natsSubscriber.Run()

}
