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

package srvhelper

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
)

type DsRefs struct {
	Blueprints []string `json:"blueprints"`
	Projects   []string `json:"projects"`
}

func toDsRefs(blueprints []*models.Blueprint) (*DsRefs, errors.Error) {
	if len(blueprints) > 0 {
		blueprintNames := make([]string, 0, len(blueprints))
		projectNames := make([]string, 0, len(blueprints))
		for _, bp := range blueprints {
			blueprintNames = append(blueprintNames, bp.Name)
			if bp.ProjectName != "" {
				projectNames = append(projectNames, bp.ProjectName)
			}
		}
		return &DsRefs{
			Blueprints: blueprintNames,
			Projects:   projectNames,
		}, errors.Conflict.New("Cannot delete record because it is referenced by blueprints")
	}
	return nil, nil
}
