package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/models"
)

const (
	AuthHeaderKey       = "X-Lake-Token"
	RefreshTimeInterval = 30 * time.Second
)

type Guard struct {
	*sync.Map
}

func NewGuard() *Guard {
	g := new(Guard)
	go func() {
		for {
			tokens, err := models.GetAuthToken()
			if err != nil && len(tokens) == 0 {
				g.Map = nil
				time.Sleep(RefreshTimeInterval)
				continue
			}
			m := new(sync.Map)
			for _, token := range tokens {
				m.Store(token.Token, struct{}{})
			}
			g.Map = m
			time.Sleep(RefreshTimeInterval)
		}
	}()
	return g
}

func (g *Guard) Auth(c *gin.Context) {
	token := c.GetHeader(AuthHeaderKey)
	if g.Map == nil {
		return
	}
	if _, ok := g.Load(token); ok {
		return
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
