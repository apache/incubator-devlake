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

package code

import (
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
)

type Note struct {
	domainlayer.DomainEntity
	PrId        string `gorm:"index;comment:References the pull request for this note;type:varchar(100)"`
	Type        string `gorm:"type:varchar(100)"`
	Author      string `gorm:"type:varchar(255)"`
	Body        string
	Resolvable  bool `gorm:"comment:Is or is not a review comment"`
	IsSystem    bool `gorm:"comment:Is or is not auto-generated vs. human generated"`
	CreatedDate time.Time
}

func (Note) TableName() string {
	return "notes"
}
