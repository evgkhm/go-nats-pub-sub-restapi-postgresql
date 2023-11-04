package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log/slog"
)

var (
	ErrPublishMsg = errors.New("NATS message did not resolve in time")
)

var Config *Conf

type Conf struct {
	url   string
	name  string
	Topic string
}

func (c Conf) Init() {
	Config = &Conf{
		url:   viper.GetString("nats.url"),
		name:  viper.GetString("nats.name"),
		Topic: viper.GetString("nats.subjects.topic"),
	}
}

func New(config *Conf, logger *slog.Logger) (*nats.Conn, jetstream.JetStream, error) {
	nc, errConnect := nats.Connect(config.url)
	if errConnect != nil {
		logger.Error("nats - New - nats.Connect:", "err", errConnect)
		return nil, nil, errConnect
	}

	js, errJetStream := jetstream.New(nc)
	if errJetStream != nil {
		logger.Error("nats - New - jetstream.New:", "err", errJetStream)
		return nil, nil, errJetStream
	}

	return nc, js, nil
}
