package nats

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
	user "go-nats-pub-sub-restapi-postgresql/gateway/internal/entity"
	"golang.org/x/net/context"
	"log"
	"sync"
)

var (
	ErrJsNil = errors.New("jetStream is nil")
)

type UserSubscribe struct {
	nc *nats.Conn
	js jetstream.JetStream
}

func NewUserSubscriber(nc *nats.Conn, js jetstream.JetStream) *UserSubscribe {
	return &UserSubscribe{
		nc: nc,
		js: js,
	}
}

type Subscriber struct {
	UserJetStream
}

type UserJetStream interface {
	Run(ctx context.Context)
	subscribe(ctx context.Context, wg *sync.WaitGroup)
	PublishMessage(js jetstream.JetStream, userDTO *user.User, topic string, message string) error
}

func NewSubscriber(nc *nats.Conn, js jetstream.JetStream) *Subscriber {
	return &Subscriber{
		UserJetStream: NewUserSubscriber(nc, js),
	}
}

func (u *UserSubscribe) subscribe(ctx context.Context, wg *sync.WaitGroup) {
	u.nc.Subscribe(Config.Topic, func(msg *nats.Msg) {
		var mqUser user.MqUser
		err := json.Unmarshal(msg.Data, &mqUser)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Consumer  =>  Subject: %s  -  ID: %d  -  Balance: %f - Method: %s\n", msg.Subject, mqUser.ID, mqUser.Balance, mqUser.Method)
	})
	wg.Done()
}

func (u *UserSubscribe) PublishMessage(js jetstream.JetStream, userDTO *user.User, topic string, message string) error {
	if js == nil {
		return fmt.Errorf("nats - PublishMessage: %w", ErrJsNil)
	}

	mqData := user.MqUser{
		ID:      userDTO.ID,
		Balance: userDTO.Balance,
		Method:  message,
	}
	b, err := json.Marshal(mqData)

	_, err = js.PublishAsync(topic, b)
	if err != nil {
		fmt.Println(err) //TODO:logger
	}

	return nil
}

func (u *UserSubscribe) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go u.subscribe(ctx, wg)
	wg.Wait()
}
