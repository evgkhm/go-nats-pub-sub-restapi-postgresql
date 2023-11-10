package transactions

import (
	"github.com/jmoiron/sqlx"
)

type Transactions struct {
	Transactor
}

//go:generate go run github.com/vektra/mockery/v2@v2.36.1 --name=Transactor
type Transactor interface {
	Rollback(tx *sqlx.Tx) error
	Commit(tx *sqlx.Tx) error
	NewTransaction() (*sqlx.Tx, error)
}

func New(postgresDB *sqlx.DB) *Transactions {
	return &Transactions{
		Transactor: NewTransaction(postgresDB),
	}
}

type Transaction struct {
	postgresDB *sqlx.DB
}

func NewTransaction(postgresDB *sqlx.DB) *Transaction {
	return &Transaction{
		postgresDB: postgresDB,
	}
}

func (t *Transaction) NewTransaction() (*sqlx.Tx, error) {
	return t.postgresDB.Beginx()
}

func (t *Transaction) Rollback(tx *sqlx.Tx) error {
	if err := tx.Rollback(); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) Commit(tx *sqlx.Tx) error {
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
