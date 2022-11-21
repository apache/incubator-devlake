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

package archived

import "github.com/apache/incubator-devlake/models/migrationscripts/archived"

type ZentaoDepartment struct {
	ConnectionId uint64  `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID           uint64  `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL" `
	Name        string `json:"name" gorm:"type:varchar(100);index"`
	Parent      int    `json:"parent" gorm:"type:varchar(100)"`
	Path        string `json:"path" gorm:"type:varchar(100)"`
	Grade       int    `json:"grade"`
	Order       int    `json:"order"`
	Position    string `json:"position" gorm:"type:varchar(100)"`
	Function    string `json:"function" gorm:"type:varchar(100)"`
	Manager     string `json:"manager" gorm:"type:varchar(100)"`
	ManagerName string `json:"managerName" gorm:"type:varchar(100)"`
	archived.NoPKModel
}
