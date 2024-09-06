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

package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addIsChildToCicdPipeline)(nil)

type cicdPipeline20240906 struct {
	IsChild bool
}

func (cicdPipeline20240906) TableName() string {
	return "cicd_pipelines"
}

type addIsChildToCicdPipeline struct{}

func (*addIsChildToCicdPipeline) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.AutoMigrate(&cicdPipeline20240906{}); err != nil {
		return err
	}
	return nil
}

func (*addIsChildToCicdPipeline) Version() uint64 {
	return 20240906120000
}

func (*addIsChildToCicdPipeline) Name() string {
	return "add is_child to cicd_pipelines"
}
