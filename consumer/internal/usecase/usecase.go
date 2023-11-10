package usecase

import (
	"context"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/repository/postgres"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/transactions"
	"log/slog"
)

type UseCase struct {
	User
}

//go:generate go run github.com/vektra/mockery/v2@v2.36.1 --name=User
type User interface {
	CreateUser(ctx context.Context, userDTO *user.User) error
	GetBalance(ctx context.Context, id string) (user.User, error)
	AccrualBalanceUser(ctx context.Context, userDTO *user.User) error
	CheckNegativeBalance(ctx context.Context, userDTO *user.User) error
	CalcNewBalance(ctx context.Context, userDTO *user.User, balance float32) float32
}

func New(repo *postgres.Repository, txService *transactions.Transactions, logger *slog.Logger) *UseCase {
	return &UseCase{
		User: NewUserUseCase(repo.UserRepository, txService.Transactor, logger),
	}
}
