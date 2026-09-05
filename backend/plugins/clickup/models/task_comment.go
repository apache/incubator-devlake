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

package models

import (
	"time"

	"github.com/apache/devlake/core/models/common"
)

// ClickUpTaskComment is the tool-layer representation of a comment on a ClickUp
// task.
//
// TODO(clickup): the comment collector/extractor is not yet implemented. This
// model exists so the table is created up-front and the convertor to
// ticket.IssueComment can be added without a follow-up migration. See
// GET /task/{id}/comment.
type ClickUpTaskComment struct {
	ConnectionId uint64     `gorm:"primaryKey"`
	Id           string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	TaskId       string     `gorm:"index;type:varchar(255)" json:"taskId"`
	Body         string     `json:"body"`
	UserId       string     `gorm:"type:varchar(255)" json:"userId"`
	CreatedDate  *time.Time `json:"createdDate"`
	common.NoPKModel
}

func (ClickUpTaskComment) TableName() string {
	return "_tool_clickup_task_comments"
}
