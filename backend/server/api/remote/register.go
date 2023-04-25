/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package remote

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/swaggo/swag"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/apache/incubator-devlake/server/services/remote"
	"github.com/apache/incubator-devlake/server/services/remote/models"
)

var (
	vld = validator.New()
)

type ApiResource struct {
	PluginName string
	Resources  map[string]map[string]plugin.ApiResourceHandler
}

// TODO add swagger doc
func RegisterPlugin(router *gin.Engine, registerEndpoints func(r *gin.Engine, pluginName string, apiResources map[string]map[string]plugin.ApiResourceHandler)) func(*gin.Context) {
	return func(c *gin.Context) {
		var pluginInfo models.PluginInfo
		if err := c.ShouldBindJSON(&pluginInfo); err != nil {
			shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
			return
		}
		if err := vld.Struct(&pluginInfo); err != nil {
			shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
			return
		}
		remotePlugin, err := remote.NewRemotePlugin(&pluginInfo)
		if err != nil {
			shared.ApiOutputError(c, errors.Default.Wrap(err, fmt.Sprintf("plugin %s could not be initialized", pluginInfo.Name)))
			return
		}
		resource := ApiResource{
			PluginName: pluginInfo.Name,
			Resources:  remotePlugin.ApiResources(),
		}
		registerEndpoints(router, resource.PluginName, resource.Resources)
		err = registerPluginOpenApiSpec(router, pluginInfo.Name, remotePlugin)
		if err != nil {
			shared.ApiOutputError(c, err)
			return
		}
		shared.ApiOutputSuccess(c, nil, http.StatusOK)
	}
}

// This function registers the open API spec provided by a plugin that implements PluginOpenApiSpec interface
// This makes make the plugin's API doc available at /plugins/swagger/<plugin-name>/index.html via swagger UI.
func registerPluginOpenApiSpec(router *gin.Engine, pluginName string, pluginOpenApiSpec plugin.PluginOpenApiSpec) errors.Error {
	spec := &swag.Spec{
		Version:          "",
		Host:             "",
		BasePath:         "",
		Schemes:          nil,
		Title:            "",
		Description:      "",
		InfoInstanceName: pluginName,
		SwaggerTemplate:  pluginOpenApiSpec.OpenApiSpec(),
	}
	swag.Register(pluginName, spec)
	router.GET(
		fmt.Sprintf("/plugins/swagger/%s/*any", pluginName),
		ginSwagger.CustomWrapHandler(
			&ginSwagger.Config{
				URL:                      "doc.json",
				DocExpansion:             "list",
				InstanceName:             pluginName,
				Title:                    fmt.Sprintf("%s API", pluginName),
				DefaultModelsExpandDepth: 1,
				DeepLinking:              true,
				PersistAuthorization:     false,
			},
			swaggerFiles.Handler,
		),
	)
	return nil
}
