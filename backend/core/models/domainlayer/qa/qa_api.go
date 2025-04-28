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

package qa

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/domainlayer"
)

// QaApi represents a QA API in the domain layer
type QaApi struct {
	domainlayer.DomainEntityExtended
	Name        string    `gorm:"type:varchar(255);comment:API name"`
	Path        string    `gorm:"type:varchar(255);comment:API path"`
	Method      string    `gorm:"type:varchar(255);comment:API method"`
	CreateTime  time.Time `gorm:"comment:API creation time"`
	CreatorId   string    `gorm:"type:varchar(255);comment:Creator ID"`
	QaProjectId string    `gorm:"type:varchar(255);index;comment:Project ID"`
}

func (QaApi) TableName() string {
	return "qa_apis"
}
