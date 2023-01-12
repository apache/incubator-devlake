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

package runner

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	contextimpl "github.com/apache/incubator-devlake/impls/context"
	"github.com/apache/incubator-devlake/impls/dalgorm"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"gorm.io/gorm"
	"sync"
)

var app_lock sync.Mutex
var app_inited bool

// CreateAppBasicRes returns a application level BasicRes instance based on .env/environment variables
// it is useful because multiple places need BasicRes including `main.go` `directrun` and `worker`
// keep in mind this function can be called only once
func CreateAppBasicRes() context.BasicRes {
	app_lock.Lock()
	if app_inited {
		panic(fmt.Errorf("CreateAppBasicRes can be called once"))
	}
	app_inited = true
	app_lock.Unlock()
	cfg := config.GetConfig()
	logger := logruslog.Global
	db, err := NewGormDb(cfg, logger)
	if err != nil {
		panic(err)
	}
	dalgorm.Init(cfg.GetString(plugin.EncodeKeyEnvStr))
	return CreateBasicRes(cfg, logger, db)
}

// CreateBasicRes returns a BasicRes based on what was given
func CreateBasicRes(cfg config.ConfigReader, logger log.Logger, db *gorm.DB) context.BasicRes {
	return contextimpl.NewDefaultBasicRes(cfg, logger, dalgorm.NewDalgorm(db))
}
