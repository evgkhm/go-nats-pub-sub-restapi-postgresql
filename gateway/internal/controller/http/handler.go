package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go/jetstream"
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/controller/mq/nats"
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

func (h Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/create_user", h.createUser)
	router.GET("/get_balance_user/:id", h.getBalanceUserByID)
	router.POST("/accrual_balance_user", h.accrualBalanceUser)

	return router
}

func New(js jetstream.JetStream, natsSubscriber *nats.Subscriber, logger *slog.Logger) *Handler {
	return &Handler{
		js:             js,
		natsSubscriber: natsSubscriber,
		logger:         logger,
	}
}
