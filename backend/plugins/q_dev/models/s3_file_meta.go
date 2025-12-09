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

	"github.com/apache/incubator-devlake/core/models/common"
)

// QDevS3FileMeta 存储S3文件的元数据信息
type QDevS3FileMeta struct {
	common.NoPKModel
	ConnectionId  uint64     `gorm:"primaryKey"`
	FileName      string     `gorm:"primaryKey;type:varchar(255)"`
	S3Path        string     `gorm:"type:varchar(512)" json:"s3Path"`
	ScopeId       string     `gorm:"type:varchar(255);index" json:"scopeId"`
	Processed     bool       `gorm:"default:false"`
	ProcessedTime *time.Time `gorm:"default:null"`
}

func (QDevS3FileMeta) TableName() string {
	return "_tool_q_dev_s3_file_meta"
}
