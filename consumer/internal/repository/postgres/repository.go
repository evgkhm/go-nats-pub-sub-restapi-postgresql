package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
)

type Repository struct {
	UserRepository
}

//go:generate go run github.com/vektra/mockery/v2@v2.36.1 --name=UserRepository
type UserRepository interface {
	CreateUser(ctx context.Context, userDTO *user.User) error
	GetBalance(ctx context.Context, id uint64) (float32, error)
	UserBalanceAccrual(ctx context.Context, userDTO *user.User) error
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(db),
	}
}
