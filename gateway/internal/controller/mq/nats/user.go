package nats

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
	user "go-nats-pub-sub-restapi-postgresql/gateway/internal/entity"
	"golang.org/x/net/context"
	"log/slog"
	"sync"
)

var (
	ErrJsNil = errors.New("jetStream is nil")
)

type UserSubscribe struct {
	nc     *nats.Conn
	js     jetstream.JetStream
	logger *slog.Logger
}

func NewUserSubscriber(nc *nats.Conn, js jetstream.JetStream, logger *slog.Logger) *UserSubscribe {
	return &UserSubscribe{
		nc:     nc,
		js:     js,
		logger: logger,
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

func NewSubscriber(nc *nats.Conn, js jetstream.JetStream, logger *slog.Logger) *Subscriber {
	return &Subscriber{
		UserJetStream: NewUserSubscriber(nc, js, logger),
	}
}

func (u *UserSubscribe) subscribe(ctx context.Context, wg *sync.WaitGroup) {
	u.nc.Subscribe(Config.Topic, func(msg *nats.Msg) {
		var mqUser user.MqUser
		err := json.Unmarshal(msg.Data, &mqUser)
		if err != nil {
			u.logger.Error("nats - UserSubscribe - subscribe - json.Unmarshal:", "err", err)
		}
		u.logger.Info("Consumer",
			slog.String("Subject", msg.Subject),
			slog.Uint64("ID", mqUser.ID),
			slog.Float64("Balance", float64(mqUser.Balance)),
			slog.String("Method", mqUser.Method))
	})
	wg.Done()
}

func (u *UserSubscribe) PublishMessage(js jetstream.JetStream, userDTO *user.User, topic string, message string) error {
	if js == nil {
		u.logger.Error("nats - UserSubscribe - PublishMessage:", "err", ErrJsNil)
		return ErrJsNil
	}

	mqData := user.MqUser{
		ID:      userDTO.ID,
		Balance: userDTO.Balance,
		Method:  message,
	}
	b, err := json.Marshal(mqData)
	if err != nil {
		u.logger.Error("nats - UserSubscribe - PublishMessage - json.Marshal:", "err", err)
		return err
	}

	_, err = js.PublishAsync(topic, b)
	if err != nil {
		u.logger.Error("nats - UserSubscribe - PublishMessage - js.PublishAsync:", "err", err)
	}

	return nil
}

func (u *UserSubscribe) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go u.subscribe(ctx, wg)
	wg.Wait()
}
