package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go/jetstream"
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/controller/mq/nats"
	"go-nats-pub-sub-restapi-postgresql/gateway/pkg/validator"
	"log/slog"
)

type Creating interface {
	CreateUser(c *gin.Context)
}

type Handler struct {
	js             jetstream.JetStream
	natsSubscriber *nats.Subscriber
	logger         *slog.Logger
}

func (h *Handler) InitRoutes() (*gin.Engine, error) {
	router := gin.Default()

	err := validator.New()
	if err != nil {
		return nil, err
	}

	router.POST("/create_user", h.createUser)
	router.GET("/get_balance_user/:id", h.getBalanceUserByID)
	router.POST("/accrual_balance_user", h.accrualBalanceUser)

	return router, nil
}

func New(js jetstream.JetStream, natsSubscriber *nats.Subscriber, logger *slog.Logger) *Handler {
	return &Handler{
		js:             js,
		natsSubscriber: natsSubscriber,
		logger:         logger,
	}
}

func (h *Handler) Validate(vld validator.Validator) error {
	if err := vld.Validate(); err != nil {
		vErr := errors.New("validation error")

		var vldErr validator.ValidationError

		if errors.As(err, &vldErr) {
			for _, v := range vldErr.Fields {
				vErr = errors.Join(vErr, errors.New(v))
			}
		}

		return vErr
	}

	return nil
}
