package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/models"
)

const (
	AuthHeaderKey = "X-Lake-Token"
)

func Auth(c *gin.Context) {
	token := c.GetHeader(AuthHeaderKey)
	if checkToken(token) {
		return
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func checkToken(token string) bool {
	authToken, err := models.GetAuthToken(token)
	if err != nil || authToken == nil {
		return false
	}
	return true
}
