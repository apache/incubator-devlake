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

package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/server/services/remote/models/migrationscripts/azuredevops"
)

var allMigrations = map[string][]plugin.MigrationScript{
	"azuredevops": {
		&azuredevops.AddRawDataForScope{},
		&azuredevops.DecryptConnectionFields{},
	},
}

// All return Go-defined migration scripts of remote plugins. These migrations are intended for more advanced
// use cases where they need to be defined in the Go language. Example: a migration that relies on checking existing data
func All(pluginName string) []plugin.MigrationScript {
	if all, ok := allMigrations[pluginName]; ok {
		return all
	}
	return nil
}
