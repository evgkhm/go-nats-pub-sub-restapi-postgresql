package nats

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/usecase"
	"golang.org/x/net/context"
	"log"
	"runtime"
	"strconv"
	"sync"
)

var (
	ErrJsNil = errors.New("jetStream is nil")
)

type UserSubscribe struct {
	nc      *nats.Conn
	js      jetstream.JetStream
	useCase *usecase.UseCase
}

type Subscriber struct {
	UserJetStream
}

type UserJetStream interface {
	Run()
	createUser(ctx context.Context, mqUser user.MqUser)
	subscribe(wg *sync.WaitGroup)
	publishMessage(ctx context.Context, userDTO *user.User, message string) error
	getBalanceUser(ctx context.Context, mqUser user.MqUser)
	accrualBalanceUser(ctx context.Context, mqUser user.MqUser)
}

func NewSubscriber(nc *nats.Conn, js jetstream.JetStream, useCase *usecase.UseCase) *Subscriber {
	return &Subscriber{
		UserJetStream: NewUserSubscriber(nc, js, useCase),
	}
}

func NewUserSubscriber(nc *nats.Conn, js jetstream.JetStream, useCase *usecase.UseCase) *UserSubscribe {
	return &UserSubscribe{
		nc:      nc,
		js:      js,
		useCase: useCase,
	}
}

func (u *UserSubscribe) publishMessage(ctx context.Context, userDTO *user.User, message string) error {
	if u.js == nil {
		return fmt.Errorf("nats - PublishMessage: %w", ErrJsNil)
	}

	mqData := user.MqUser{
		ID:      userDTO.ID,
		Balance: userDTO.Balance,
		Method:  message,
	}
	b, err := json.Marshal(mqData)
	if err != nil {
		fmt.Println(err) //TODO:logger
	}

	_, err = u.js.PublishAsync(Config.Topic, b)
	if err != nil {
		fmt.Println(err) //TODO:logger
	}

	return nil
}

func (u *UserSubscribe) createUser(ctx context.Context, mqUser user.MqUser) { //context
	userDTO := &user.User{
		ID:      mqUser.ID,
		Balance: mqUser.Balance,
	}

	err := u.useCase.CreateUser(ctx, userDTO)
	if err != nil {
		u.publishMessage(ctx, userDTO, err.Error())
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

func (u *UserSubscribe) subscribe(wg *sync.WaitGroup) {
	u.nc.Subscribe(Config.Topic, func(msg *nats.Msg) {
		var mqUser user.MqUser
		err := json.Unmarshal(msg.Data, &mqUser)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Consumer  =>  Subject: %s  -  ID: %d  -  Balance: %f - Method: %s\n", msg.Subject, mqUser.ID, mqUser.Balance, mqUser.Method)

		ctx := context.Background()
		switch mqUser.Method {
		case "Create user":
			u.createUser(ctx, mqUser)
		case "Get user balance":
			u.getBalanceUser(ctx, mqUser)
		case "Accrual user balance":
			u.accrualBalanceUser(ctx, mqUser)
		}
	})
	wg.Done()
}

func (u *UserSubscribe) Run() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go u.subscribe(wg)
	wg.Wait()
	runtime.Goexit()
}
