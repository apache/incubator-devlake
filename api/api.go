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

package api

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/apache/incubator-devlake/api/docs"
	"github.com/apache/incubator-devlake/api/shared"
	"github.com/apache/incubator-devlake/logger"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title  DevLake Swagger API
// @version 0.1
// @description  <h2>This is the main page of devlake api</h2>
// sdfasdfasd
// @license.name Apache-2.0
// @host localhost:8080
// @BasePath /
func CreateApiService() {
	services.Init()
	v := config.GetConfig()
	gin.SetMode(v.GetString("MODE"))
	router := gin.Default()

	// Wait for user confirmation if db migration is needed
	router.GET("/proceed-db-migration", func(ctx *gin.Context) {
		if !services.MigrationRequireConfirmation() {
			shared.ApiOutputError(ctx, fmt.Errorf("no pending migration"), http.StatusBadRequest)
			return
		}
		err := services.ExecuteMigration()
		if err != nil {
			shared.ApiOutputError(ctx, err, http.StatusBadRequest)
			return
		}
		shared.ApiOutputSuccess(ctx, nil, http.StatusOK)
	})
	router.Use(func(ctx *gin.Context) {
		if !services.MigrationRequireConfirmation() {
			return
		}
		shared.ApiOutputError(
			ctx,
			fmt.Errorf("Database migration is required for Apache DevLake to function properly, it might cause the "+
				"collected data gets wiped out for consistency. Please send a request to `/proceed-migrations` "+
				"if it is ok, or you may downgrade back to the older version you previous used"),
			http.StatusPreconditionRequired,
		)
		ctx.Abort()
		return
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//endpoint debug log
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.Global.Printf("endpoint %v %v %v %v", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	// CORS CONFIG
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           120 * time.Hour,
	}))

	RegisterRouter(router)
	err := router.Run(v.GetString("PORT"))
	if err != nil {
		panic(err)
	}
}
