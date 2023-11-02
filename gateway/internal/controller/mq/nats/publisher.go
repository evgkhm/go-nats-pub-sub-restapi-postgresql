package nats

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
	user "go-nats-pub-sub-restapi-postgresql/gateway/internal/entity"
)

var (
	ErrJsNil = errors.New("jetStream is nil")
)

func PublishMessage(js jetstream.JetStream, userDTO *user.User, topic string, message string) error {
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
