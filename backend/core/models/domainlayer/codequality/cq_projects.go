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

package codequality

import (
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.Scope = (*CqProject)(nil)

type CqProject struct {
	domainlayer.DomainEntityExtended
	Name             string `gorm:"type:varchar(255)"`
	Qualifier        string `gorm:"type:varchar(255)"`
	Visibility       string `gorm:"type:varchar(64)"`
	LastAnalysisDate *common.Iso8601Time
	CommitSha        string `gorm:"type:varchar(128)"`
}

func (CqProject) TableName() string {
	return "cq_projects"
}

func (s *CqProject) ScopeId() string {
	return s.Id
}

func (s *CqProject) ScopeName() string {
	return s.Name
}
