package transactions

import (
	"github.com/jmoiron/sqlx"
)

type TransactionServiceImpl struct {
	postgresDB *sqlx.DB
}

func New(postgresDB *sqlx.DB) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		postgresDB: postgresDB,
	}
}

func (t *TransactionServiceImpl) NewTransaction() (*sqlx.Tx, error) {
	return t.postgresDB.Beginx()
}

func (t *TransactionServiceImpl) Rollback(tx *sqlx.Tx) error {
	if err := tx.Rollback(); err != nil {
		return err
	}
	return nil
}

func (t TransactionServiceImpl) Commit(tx *sqlx.Tx) error {
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
