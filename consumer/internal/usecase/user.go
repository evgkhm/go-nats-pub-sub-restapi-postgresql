package usecase

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/repository/postgres"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/transactions"
	"log/slog"
	"strconv"
)

var (
	ErrUserAccrualNegativeBalance = errors.New("you cannot accrual with negative balance")
)

type UserUseCase struct {
	userRepo  postgres.UserRepository
	txService transactions.Transactor
	logger    *slog.Logger
}

func NewUserUseCase(userRepo postgres.UserRepository, txService transactions.Transactor, logger *slog.Logger) *UserUseCase {
	return &UserUseCase{
		userRepo:  userRepo,
		txService: txService,
		logger:    logger,
	}
}

func (u *UserUseCase) AccrualBalanceUser(ctx context.Context, userDTO *user.User) error {
	tx, err := u.txService.NewTransaction()
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			u.logger.Error("usecase - UseCase - UserBalanceAccrual - u.txService.NewTransaction - u.txService.Rollback", "err", err)
		}
		u.logger.Error("usecase - UseCase - UserBalanceAccrual - u.txService.NewTransaction", "err", err)
		return err
	}

	// Проверка на отрицательный баланс
	errNegBlnc := u.CheckNegativeBalance(ctx, userDTO)
	if errNegBlnc != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			u.logger.Error("usecase - UseCase - UserBalanceAccrual - u.txService.Rollback", "err", err)
			return fmt.Errorf("usecase - UseCase - UserBalanceAccrual - u.txService.Rollback: %w", err)
		}
		u.logger.Error("usecase - UseCase - UserBalanceAccrual", "err", ErrUserAccrualNegativeBalance)
		return ErrUserAccrualNegativeBalance
	}

	id := userDTO.ID
	// Узнать текущий баланс
	currBalance, err := u.userRepo.GetBalance(ctx, strconv.FormatUint(id, 10), tx)
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			u.logger.Error("usecase - UseCase - UserBalanceAccrual - u.repo.GetBalance - u.txService.Rollback", "err", err)
		}
		u.logger.Error("usecase - UseCase - UserBalanceAccrual - u.repo.GetBalance", "err", err)
		return err
	}

	// Начислить новый баланс
	userDTO.Balance = u.CalcNewBalance(ctx, userDTO, currBalance)
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			u.logger.Error("usecase - UseCase - UserBalanceAccrual - calcNewBalance - u.txService.Rollback", "err", err)
		}
		u.logger.Error("usecase - UseCase - UserBalanceAccrual - calcNewBalance", "err", err)
		return err
	}

	err = u.userRepo.UserBalanceAccrual(ctx, tx, userDTO)
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			u.logger.Error("usecase - UseCase - UserBalanceAccrual - u.repo.UserBalanceAccrual - u.txService.Rollback", "err", err)
			return err
		}
		u.logger.Error("usecase - UseCase - UserBalanceAccrual - u.repo.UserBalanceAccrual", "err", err)
	}

	return u.txService.Commit(tx)
}

func (u *UserUseCase) CalcNewBalance(ctx context.Context, userDTO *user.User, balance float32) float32 {
	newBalance := userDTO.Balance + balance
	return newBalance
}

func (u *UserUseCase) CreateUser(ctx context.Context, userDTO *user.User) error {
	tx, err := u.txService.NewTransaction()
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			u.logger.Error("usecase - UseCase - CreateUser - u.txService.NewTransaction - u.txService.Rollback", "err", err)
			return err
		}
		u.logger.Error("usecase - UseCase - CreateUser - u.txService.NewTransaction", "err", err)
		return err
	}

	err = u.userRepo.CreateUser(ctx, tx, userDTO)
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			u.logger.Error("usecase - UseCase - CreateUser - u.repo.CreateUser - u.txService.Rollback", "err", err)
			return err
		}
		u.logger.Error("usecase - UseCase - CreateUser - u.repo.CreateUser", "err", err)
		return err
	}

	return u.txService.Commit(tx)
}

func (u *UserUseCase) GetBalance(ctx context.Context, id string) (user.User, error) {
	var userDTO user.User
	tx, err := u.txService.NewTransaction()
	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			u.logger.Error("usecase - UseCase - GetBalance - u.txService.NewTransaction - u.txService.Rollback", "err", err)
			return user.User{}, err
		}
		u.logger.Error("usecase - UseCase - GetBalance - u.txService.NewTransaction", "err", err)
		return user.User{}, err
	}

	balance, err := u.userRepo.GetBalance(ctx, id, tx)

	if err != nil {
		errRollback := u.txService.Rollback(tx)
		if errRollback != nil {
			u.logger.Error("usecase - UseCase - GetBalance - u.repo.GetBalance - u.txService.Rollback", "err", err)
			return user.User{}, err
		}
		u.logger.Error("usecase - UseCase - GetBalance - u.repo.GetBalance", "err", err)
		return user.User{}, err
	}

	userDTO.Balance = balance
	userDTO.ID, _ = strconv.ParseUint(id, 10, 64)

	u.txService.Commit(tx)

	return userDTO, nil
}

func (u *UserUseCase) CheckNegativeBalance(ctx context.Context, userDTO *user.User) error {
	if userDTO.Balance < 0 {
		u.logger.Error("usecase - UseCase - UserBalanceAccrual", ErrUserAccrualNegativeBalance)
		return ErrUserAccrualNegativeBalance
	}
	return nil
}
