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

package azuredevops

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*DecryptConnectionFields)(nil)

type azureDevopsConnection20240201 struct {
	archived.Model
	Name string `gorm:"size:255;uniqueIndex"`
}

type UniqueConnectionName struct{}

func (script *UniqueConnectionName) Up(basicRes context.BasicRes) errors.Error {
	connectionCounter := 1
	err := migrationhelper.TransformColumns(
		basicRes,
		script,
		"_tool_azuredevops_azuredevopsconnections",
		[]string{"name"},
		func(src *azureDevopsConnection20240201) (*azureDevopsConnection20240201, errors.Error) {
			src.Name = fmt.Sprintf("%s_%d", src.Name, connectionCounter)
			connectionCounter++
			return src, nil
		},
	)
	return err
}

func (*UniqueConnectionName) Version() uint64 {
	return 20240201000001
}

func (script *UniqueConnectionName) Name() string {
	return "Make AzDo connection name unique"
}
