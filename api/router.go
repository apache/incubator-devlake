package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/api/env"
	"github.com/merico-dev/lake/api/task"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/services"
)

func RegisterRouter(r *gin.Engine) {
	r.POST("/task", task.Post)
	r.GET("/task", task.Get)
	r.GET("/task/pending", task.GetPending)
	r.DELETE("/task/:taskId", task.Delete)
	r.POST("/env", env.Set)
	r.GET("/env", env.Get)

	// mount all api resources for all plugins
	pluginsApiResources, err := services.GetPluginsApiResources()
	if err != nil {
		panic(err)
	}
	for pluginName, apiResources := range pluginsApiResources {
		for resourcePath, resourceHandlers := range apiResources {
			for method, h := range resourceHandlers {
				handler := h // block scoping
				r.Handle(
					method,
					fmt.Sprintf("/plugins/%s/%s", pluginName, resourcePath),
					func(c *gin.Context) {
						// connect http request to plugin interface
						input := &core.ApiResourceInput{}
						if len(c.Params) > 0 {
							input.Params = make(map[string]string)
							for _, param := range c.Params {
								input.Params[param.Key] = param.Value
							}
						}
						if c.Request.Body != nil {
							err := c.ShouldBindJSON(&input.Body)
							if err != nil && err.Error() != "EOF" {
								c.JSON(http.StatusBadRequest, err.Error())
								return
							}
						}
						output, err := handler(input)
						if err != nil {
							c.JSON(http.StatusBadRequest, err.Error())
						} else {
							c.JSON(http.StatusCreated, output.Body)
						}
					},
				)
			}
		}
	}
}
