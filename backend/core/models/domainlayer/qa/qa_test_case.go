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

// QaTestCase represents a QA test case in the domain layer
type QaTestCase struct {
	domainlayer.DomainEntityExtended
	Name        string    `gorm:"type:varchar(255);comment:Test case name"`
	CreateTime  time.Time `gorm:"comment:Test case creation time"`
	CreatorId   string    `gorm:"type:varchar(255);comment:Creator ID"`
	Type        string    `gorm:"type:varchar(255);comment:Test case type | functional | api"`                // enum in image, using string
	QaApiId     string    `gorm:"type:varchar(255);comment:Valid only when type = api, represents qa_api_id"` // nullable in image, using string
	QaProjectId string    `gorm:"type:varchar(255);index;comment:Project ID"`
}

func (qaTestCase *QaTestCase) TableName() string {
	return "qa_test_cases"
}
