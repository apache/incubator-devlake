package version

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/version"
)

func Get(c *gin.Context) {
	c.JSON(http.StatusOK, struct {
		Version string `json:"version"`
	}{
		Version: version.Version,
	})
}
