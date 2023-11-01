package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
)

type Repository struct {
	UserRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, tx *sqlx.Tx, userDTO *user.User) error
	GetBalance(ctx context.Context, id string, tx *sqlx.Tx) (float32, error)
	UserBalanceAccrual(ctx context.Context, tx *sqlx.Tx, userDTO *user.User) error
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(db),
	}
}
