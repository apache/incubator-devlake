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

package test

import (
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/core"
	gitlabhack "github.com/apache/incubator-devlake/plugins/gitlab/api"
	"github.com/apache/incubator-devlake/services"
	"sync"
)

var once = &sync.Once{}

// Setup sets up the environment for integration/e2e non-plugin tests. Wipes the DB and reruns all migrations.
// NOTE: all plugins tests must run before any tests that use this function, otherwise there will be duplicate migration attempts!
func Setup() {
	once.Do(func() {
		v := config.GetConfig()
		encKey := v.GetString(core.EncodeKeyEnvStr)
		if encKey == "" {
			// Randomly generate a bunch of encryption keys and set them to config
			encKey = core.RandomEncKey()
			v.Set(core.EncodeKeyEnvStr, encKey)
			err := config.WriteConfig(v)
			if err != nil {
				panic(err)
			}
		}
		tester := e2ehelper.NewDataFlowTester(nil, "setup", nil)
		// cleanup any existing data and migrations (e.g. plugins manual migrations)
		tester.DropAllTables()
		// run migrations
		gitlabhack.Init(config.GetConfig(), tester.Log, tester.Db)
		services.Init()
	})
}
