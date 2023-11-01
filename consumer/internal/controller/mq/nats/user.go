package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/usecase"
	"golang.org/x/net/context"
	"strconv"
)

//type User interface {
//	CreateUser()
//	GetBalance()
//	AccrualBalanceUser()
//}

type UserSubscribe struct {
	nc      *nats.Conn
	js      jetstream.JetStream
	useCase *usecase.UseCase
}

func NewUserSubscriber(nc *nats.Conn, js jetstream.JetStream, useCase *usecase.UseCase) *UserSubscribe {
	return &UserSubscribe{
		nc:      nc,
		js:      js,
		useCase: useCase,
	}
}

func (u *UserSubscribe) createUser(ctx context.Context, mqUser user.MqUser) { //context
	userDTO := &user.User{
		ID:      mqUser.ID,
		Balance: mqUser.Balance,
	}

	err := u.useCase.CreateUser(ctx, userDTO)
	if err != nil {
		u.publishMessage(ctx, userDTO, err.Error())
		//nats.PublishMessage(h.js, userDTO, nats2.Config.Topic, ErrInternalServer.Error())
		return
	}

	err = u.publishMessage(ctx, userDTO, "User created")
	if err != nil {
		return
	}
}

func (u *UserSubscribe) getBalanceUser(ctx context.Context, mqUser user.MqUser) { //context
	id := mqUser.ID

	userDTO, err := u.useCase.GetBalance(ctx, strconv.FormatUint(id, 10))
	if err != nil {
		u.publishMessage(ctx, &userDTO, err.Error())
		return
	}

	err = u.publishMessage(ctx, &userDTO, "Got user balance")
	if err != nil {
		//TODO:log
		return
	}
}

func (u *UserSubscribe) accrualBalanceUser(ctx context.Context, mqUser user.MqUser) {
	userDTO := &user.User{
		ID:      mqUser.ID,
		Balance: mqUser.Balance,
	}

	err := u.useCase.AccrualBalanceUser(ctx, userDTO)
	if err != nil {
		u.publishMessage(ctx, userDTO, err.Error())
		return
	}

	err = u.publishMessage(ctx, userDTO, "Accrual user balance done")
	if err != nil {
		//TODO:log
		return
	}
}
