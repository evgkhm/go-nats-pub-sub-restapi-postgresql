package usecase

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/repository/postgres"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/transactions"
	"strconv"
)

var (
	ErrUserAccrualNegativeBalance = errors.New("you cannot accrual with negative balance")
)

type UserUseCase struct {
	userRepo  postgres.UserRepository
	txService *transactions.TransactionServiceImpl
}

func NewUserUseCase(userRepo postgres.UserRepository, txService *transactions.TransactionServiceImpl) *UserUseCase {
	return &UserUseCase{
		userRepo:  userRepo,
		txService: txService,
	}
}

func (u UserUseCase) AccrualBalanceUser(ctx context.Context, userDTO *user.User) error {
	tx, err := u.txService.NewTransaction()
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			return fmt.Errorf("usecase - UseCase - UserBalanceAccrual - u.txService.NewTransaction - u.txService.Rollback: %w", err)
		}
		return fmt.Errorf("usecase - UseCase - UserBalanceAccrual - u.txService.NewTransaction: %w", err)
	}

	// Проверка на отрицательный баланс
	if userDTO.Balance < 0 {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			return fmt.Errorf("usecase - UseCase - UserBalanceAccrual - u.txService.Rollback: %w", err)
		}
		return fmt.Errorf("usecase - UseCase - UserBalanceAccrual: %w", ErrUserAccrualNegativeBalance)
	}

	id := userDTO.ID
	// Узнать текущий баланс
	currBalance, err := u.userRepo.GetBalance(ctx, strconv.FormatUint(id, 10), tx)
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			return fmt.Errorf("usecase - UseCase - UserBalanceAccrual - u.repo.GetBalance - u.txService.Rollback: %w", err)
		}
		return fmt.Errorf("usecase - UseCase - UserBalanceAccrual - u.repo.GetBalance: %w", err)
	}

	// Начислить новый баланс
	newBalance := userDTO.Balance + currBalance
	userDTO.Balance = newBalance

	err = u.userRepo.UserBalanceAccrual(ctx, tx, userDTO)
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			return fmt.Errorf("usecase - UseCase - UserBalanceAccrual - u.repo.UserBalanceAccrual - u.txService.Rollback: %w", err)
		}
		return fmt.Errorf("usecase - UseCase - UserBalanceAccrual - u.repo.UserBalanceAccrual: %w", err)
	}

	return u.txService.Commit(tx)
}

func (u UserUseCase) CreateUser(ctx context.Context, userDTO *user.User) error {
	tx, err := u.txService.NewTransaction()
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			return fmt.Errorf("usecase - UseCase - CreateUser - u.txService.NewTransaction - u.txService.Rollback: %w", err)
		}
		return fmt.Errorf("usecase - UseCase - CreateUser - u.txService.NewTransaction: %w", err)
	}

	err = u.userRepo.CreateUser(ctx, tx, userDTO)
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			return fmt.Errorf("usecase - UseCase - CreateUser - u.repo.CreateUser - u.txService.Rollback: %w", err)
		}
		return fmt.Errorf("usecase - UseCase - CreateUser - u.txService.CreateUser: %w", err)
	}

	return u.txService.Commit(tx)
}

func (u UserUseCase) GetBalance(ctx context.Context, id string) (user.User, error) {
	var userDTO user.User
	tx, err := u.txService.NewTransaction()
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			return user.User{}, fmt.Errorf("usecase - UseCase - GetBalance - u.txService.NewTransaction - u.txService.Rollback: %w", err)
		}
		return user.User{}, fmt.Errorf("usecase - UseCase - GetBalance - u.txService.NewTransaction: %w", err)
	}

	balance, err := u.userRepo.GetBalance(ctx, id, tx)

	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			return user.User{}, fmt.Errorf("usecase - UseCase - GetBalance - u.repo.GetBalance - u.txService.Rollback: %w", err)
		}
		return user.User{}, fmt.Errorf("usecase - UseCase - GetBalance - u.txService.GetBalance: %w", err)
	}

	userDTO.Balance = balance
	userDTO.ID, _ = strconv.ParseUint(id, 10, 64)

	u.txService.Commit(tx)

	return userDTO, nil
}
