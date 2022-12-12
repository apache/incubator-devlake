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

package ticket

import (
	"time"

	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/core"

	"github.com/apache/incubator-devlake/models/domainlayer"
)

var _ core.Scope = (*Board)(nil)

type Board struct {
	domainlayer.DomainEntity
	Name        string `gorm:"type:varchar(255)"`
	Description string
	Url         string `gorm:"type:varchar(255)"`
	CreatedDate *time.Time
	Type        string `gorm:"type:varchar(255)"`
}

func (Board) TableName() string {
	return "boards"
}

func (r *Board) ScopeId() string {
	return r.Id
}

func (r *Board) ScopeName() string {
	return r.Name
}

func NewBoard(id string, name string) *Board {
	board := &Board{
		DomainEntity: domainlayer.NewDomainEntity(id),
	}

	board.Name = name
	board.CreatedDate = &board.CreatedAt

	return board
}

type BoardSprint struct {
	common.NoPKModel
	BoardId  string `gorm:"primaryKey;type:varchar(255)"`
	SprintId string `gorm:"primaryKey;type:varchar(255)"`
}

func (BoardSprint) TableName() string {
	return "board_sprints"
}
