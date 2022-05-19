package version

import (
	"net/http"

	"github.com/apache/incubator-devlake/version"
	"github.com/gin-gonic/gin"
)

// @Summary Get the version of lake
// @Description return a object
// @Tags version
// @Accept application/json
// @Success 200  {string} json ""
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /version [get]
func Get(c *gin.Context) {
	c.JSON(http.StatusOK, struct {
		Version string `json:"version"`
	}{
		Version: version.Version,
	})
}
