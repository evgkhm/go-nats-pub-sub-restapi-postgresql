package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go-nats-pub-sub-restapi-postgresql/gateway/internal/controller/mq/nats"
	user "go-nats-pub-sub-restapi-postgresql/gateway/internal/entity"
	"go-nats-pub-sub-restapi-postgresql/gateway/pkg/validator"
	"net/http"
	"strconv"
)

var (
	ErrInternalServer = errors.New("InternalServerError")
	ErrBadRequest     = errors.New("Bad request")
	ErrEmptyCodeParam = errors.New("Empty code param")
)

func (h *Handler) createUser(c *gin.Context) {
	var userDTO *user.UserWithBalance
	err := c.BindJSON(&userDTO)
	if err != nil {
		h.logger.Error("Error while bind userDTO to JSON", "err", ErrBadRequest)
		c.JSON(http.StatusBadRequest, ErrBadRequest.Error())
		return
	}
	userDTO.Method = "Create user"

	if err = h.Validate(validator.StructValidator(userDTO)); err != nil {
		h.logger.Error("Error while validate userDTO", "err", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.natsSubscriber.PublishMessage(h.js, userDTO, nats.Config.Topic)
	if err != nil {
		h.logger.Error("Error while publish message to nats", "err", ErrInternalServer)
		c.JSON(http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	c.JSON(http.StatusOK, "Ok")
}

func (h *Handler) getBalanceUserByID(c *gin.Context) {
	var err error
	id := c.Param("id")

	var userDTO user.User
	userDTO.ID, err = strconv.ParseUint(id, 10, 64)
	userDTO.Method = "Get user balance"

	if err = h.Validate(validator.StructValidator(userDTO)); err != nil {
		h.logger.Error("Error while validate userDTO", "err", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.natsSubscriber.PublishMessage(h.js, &userDTO, nats.Config.Topic)
	if err != nil {
		h.logger.Error("Error while publish message to nats", "err", ErrInternalServer)
		c.JSON(http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}

	c.JSON(http.StatusOK, "Ok")
}

func (h *Handler) accrualBalanceUser(c *gin.Context) {
	var userDTO *user.UserWithBalance
	err := c.BindJSON(&userDTO)
	if err != nil {
		h.logger.Error("Error while bind userDTO to JSON", "err", ErrBadRequest)
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}
	userDTO.Method = "Accrual user balance"

	if err = h.Validate(validator.StructValidator(userDTO)); err != nil {
		h.logger.Error("Error while validate userDTO", "err", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.natsSubscriber.PublishMessage(h.js, userDTO, nats.Config.Topic)
	if err != nil {
		h.logger.Error("Error while publish message to nats", "err", ErrInternalServer)
		c.JSON(http.StatusInternalServerError, ErrInternalServer.Error())
		return
	}
	c.JSON(http.StatusOK, "Ok")
}
