package nats

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
	"go-nats-pub-sub-restapi-postgresql/consumer/internal/usecase"
	"golang.org/x/net/context"
	"log"
	"runtime"
	"sync"
)

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

func (u *UserSubscribe) subscribe(wg *sync.WaitGroup) {
	//wg := &sync.WaitGroup{}
	//wg.Add(1)
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
		//wg.Done()
	})
	wg.Done()
	//s.nc.Flush()
	//wg.Wait()

}

func (u *UserSubscribe) Run() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go u.subscribe(wg)
	wg.Wait()
	runtime.Goexit()
}
