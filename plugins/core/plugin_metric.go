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

package core

import (
	"github.com/apache/incubator-devlake/errors"
)

/*
 */
type Tabler interface {
	TableName() string
}

type PluginMetric interface {
	// returns a list of required data entities and expected features.
	// [{ "model": "cicd_tasks", "requiredFields": {"column": "type", "execptedValue": "Deployment"}}, ...]
	RequiredDataEntities() (data []map[string]interface{}, err errors.Error)

	// This method returns all models of the current plugin
	GetTablesInfo() []Tabler

	// returns if the metric depends on Project for calculation.
	// Currently, only dora would return true.
	IsProjectMetric() bool

	// indicates which plugins must be executed before executing this one.
	// declare a set of dependencies with this
	RunAfter() ([]string, errors.Error)

	// returns an empty pointer of the plugin setting struct.
	// (no concrete usage at this point)
	Settings() (p interface{})
}
