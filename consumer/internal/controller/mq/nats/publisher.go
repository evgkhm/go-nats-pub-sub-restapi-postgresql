package nats

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
	"golang.org/x/net/context"
)

var (
	ErrJsNil = errors.New("jetStream is nil")
)

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
