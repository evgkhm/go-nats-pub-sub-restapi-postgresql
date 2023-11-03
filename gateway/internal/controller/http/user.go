package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/controller/mq/nats"
	user "go-nats-pub-sub-restapi-postgresql/gateway/internal/entity"
	"net/http"
	"strconv"
)

var (
	ErrInternalServer = errors.New("InternalServerError")
	ErrBadRequest     = errors.New("Bad request")
)

func (h Handler) createUser(c *gin.Context) {
	var userDTO *user.User
	err := c.BindJSON(&userDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequest.Error())
		return
	}
	err = h.natsSubscriber.PublishMessage(h.js, userDTO, nats.Config.Topic, "Create user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	c.JSON(http.StatusOK, "Ok")
}

func (h Handler) getBalanceUserByID(c *gin.Context) {
	var err error
	id := c.Param("id")
	var userDTO user.User
	userDTO.ID, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	err = h.natsSubscriber.PublishMessage(h.js, &userDTO, nats.Config.Topic, "Get user balance")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	c.JSON(http.StatusOK, "Ok")
}

func (h Handler) accrualBalanceUser(c *gin.Context) {
	var userDTO *user.User
	err := c.BindJSON(&userDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}

	err = h.natsSubscriber.PublishMessage(h.js, userDTO, nats.Config.Topic, "Accrual user balance")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}
	c.JSON(http.StatusOK, "Ok")
}
