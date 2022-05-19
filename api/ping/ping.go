package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Ping
// @Description check http status
// @Tags Ping
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /ping [get]
func Get(c *gin.Context) {
	c.Status(http.StatusOK)
}
