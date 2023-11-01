package usecase

import (
	"context"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/repository/postgres"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/transactions"
)

type UseCase struct {
	User
}

type User interface {
	CreateUser(ctx context.Context, userDTO *user.User) error
	GetBalance(ctx context.Context, id string) (user.User, error)
	AccrualBalanceUser(ctx context.Context, userDTO *user.User) error
}

func New(repo *postgres.Repository, txService *transactions.TransactionServiceImpl) *UseCase {
	return &UseCase{
		User: NewUserUseCase(repo.UserRepository, txService),
	}
}
