package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/api/env"
	"github.com/merico-dev/lake/api/pipelines"
	"github.com/merico-dev/lake/api/task"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/services"
)

func RegisterRouter(r *gin.Engine) {
	r.GET("/pipelines", pipelines.Index)
	r.GET("/pipelines/:pipelineId", pipelines.Get)
	r.POST("/pipelines", pipelines.Post)
	r.DELETE("/pipelines/:pipelineId", pipelines.Delete)
	r.GET("/pipelines/:pipelineId/tasks", task.Index)
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
						input.Query = c.Request.URL.Query()
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
							status := output.Status
							if status < http.StatusContinue {
								status = http.StatusOK
							}
							c.JSON(status, output.Body)
						}
					},
				)
			}
		}
	}
}
