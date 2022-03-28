package shared

import (
	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger"
)

type ApiBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ApiOutputError(c *gin.Context, err error, status int) {
	if e, ok := err.(*errors.Error); ok {
		c.JSON(e.Status, &ApiBody{
			Success: false,
			Message: err.Error(),
		})
	} else {
		logger.Global.Error("Server Internal Error: %w", err)
		c.JSON(status, &ApiBody{
			Success: false,
			Message: err.Error(),
		})
	}
	c.Writer.Header().Set("Content-Type", "application/json")
}

func ApiOutputSuccess(c *gin.Context, body interface{}, status int) {
	if body == nil {
		body = &ApiBody{
			Success: true,
			Message: "success",
		}
	}
	c.JSON(status, body)
}
