package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	ErrPublishMsg = errors.New("NATS message did not resolve in time")
)

var Config *Conf

type Conf struct {
	url  string
	name string
	//subjects string
	Topic string
}

func (c Conf) Init() {
	Config = &Conf{
		url:  viper.GetString("nats.url"),
		name: viper.GetString("nats.name"),
		//subjects: viper.GetString("nats.subjects"),
		Topic: viper.GetString("nats.subjects.topic"),
	}
}

func New(config *Conf) (*nats.Conn, jetstream.JetStream, error) {
	// connect to nats server
	nc, errConnect := nats.Connect(config.url) //config.url
	if errConnect != nil {
		return nil, nil, fmt.Errorf("nats - New - nats.Connect: %w", errConnect)
	}
	// create jetstream context from nats connection
	js, errJetStream := jetstream.New(nc)
	if errJetStream != nil {
		return nil, nil, fmt.Errorf("nats - New - jetstream.New: %w", errJetStream)
	}

	//nc.JetStream(nats.PublishAsyncMaxPending(256), jetstream.Pub)
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//
	//stream, _ := js.CreateStream(ctx, jetstream.StreamConfig{
	//	Name:     config.name,
	//	Subjects: []string{config.Topic},
	//})
	//
	//cons, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
	//	Name: "consumer",
	//})
	////_, errPublish := js.PublishAsync(config.Topic, []byte("NATS started"))
	////if errPublish != nil {
	////	return nil, nil, fmt.Errorf("nats - New - js.PublishAsync: %w", errPublish)
	////}
	//
	//wg := sync.WaitGroup{}
	//wg.Add(3)
	//
	//cc, _ := cons.Consume(func(msg jetstream.Msg) {
	//	msg.Ack()
	//	fmt.Println("received msg on", msg.Subject())
	//	wg.Done()
	//})
	//wg.Wait()
	//
	//cc.Stop()

	return nc, js, nil
}
