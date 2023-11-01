package nats

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	user "go-nats-pub-sub-restapi-postgresql/gateway/internal/entity"
	"golang.org/x/net/context"
	"log"
	"sync"
)

type Subscriber struct {
	nc *nats.Conn
}

func NewSubscriber(nc *nats.Conn) *Subscriber {
	return &Subscriber{
		nc: nc,
	}
}

func (s *Subscriber) Subscribe(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	s.nc.Subscribe(Config.Topic, func(msg *nats.Msg) {
		var mqUser user.MqUser
		err := json.Unmarshal(msg.Data, &mqUser)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Consumer  =>  Subject: %s  -  ID: %d  -  Balance: %f - Method: %s\n", msg.Subject, mqUser.ID, mqUser.Balance, mqUser.Method)
		wg.Done()
	})
	//s.nc.Flush()
	wg.Wait()

}

func (s *Subscriber) Run(ctx context.Context) {
	go s.Subscribe(ctx)
}
