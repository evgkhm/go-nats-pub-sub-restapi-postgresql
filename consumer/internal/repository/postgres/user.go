package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
)

var (
	ErrUserAlreadyExist = errors.New("such user already exists")
	ErrUserNotExist     = errors.New("such user does not exist")
)

type Repo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r Repo) CreateUser(ctx context.Context, user *user.User) error {
	var id int64
	var duplicateEntryError = &pq.Error{Code: "23505"}
	query := `INSERT INTO user_info (id, balance) VALUES ($1, $2) RETURNING id`
	row := r.db.QueryRowxContext(ctx, query, user.ID, user.Balance)
	//row := tx.QueryRowxContext(ctx, query, user.ID, user.Balance)
	err := row.Scan(&id)
	switch {
	case errors.As(err, &duplicateEntryError):
		return ErrUserAlreadyExist
	case err != nil:
		return fmt.Errorf("postgres - UsersRepositoryImpl - CreateUser - tx.QueryRowxContext - row.Scan: %w", err)
	}
	return nil
}

func (r Repo) GetBalance(ctx context.Context, id uint64) (float32, error) {
	var balance float32
	query := `SELECT balance FROM user_info WHERE id=$1 `
	err := r.db.GetContext(ctx, &balance, query, id)
	//err := tx.GetContext(ctx, &balance, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotExist // return 0, fmt.Errorf("postgres - UsersRepositoryImpl - GetBalance - tx.GetContext: %w", ErrUserNotExist)
		}
		return 0, fmt.Errorf("postgres - UsersRepositoryImpl - GetBalance - tx.GetContext: %w", err)
	}
	return balance, nil
}

func (r Repo) UserBalanceAccrual(ctx context.Context, user *user.User) error {
	query := `UPDATE user_info SET balance=$1 WHERE id=$2`
	_, err := r.db.ExecContext(ctx, query, user.Balance, user.ID)
	//_, err := tx.ExecContext(ctx, query, user.Balance, user.ID)
	if err != nil {
		return fmt.Errorf("postgres - UsersRepositoryImpl - UserBalanceAccrual - tx.ExecContext: %w", err)
	}
	return nil
}

func (r Repo) MinusBalance(ctx context.Context, user *user.User) error {
	query := `UPDATE user_info SET balance=$1 WHERE id=$2`
	_, err := r.db.ExecContext(ctx, query, user.Balance, user.ID)
	//_, err := tx.ExecContext(ctx, query, user.Balance, user.ID)
	if err != nil {
		return fmt.Errorf("postgres - ReservationRepositoryImpl - MinusBalance - tx.ExecContext: %w", err)
	}
	return nil
}
