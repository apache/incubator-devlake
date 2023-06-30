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
	"fmt"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addRawParamTableForScopes)(nil)

type addRawParamTableForScopes struct {
	p plugin.PluginMeta
}

func newAddRawParamTableForScopes(p plugin.PluginMeta) plugin.MigrationScript {
	if src, ok := p.(plugin.PluginSource); ok {
		if scope := src.Scope(); scope != nil {
			return &addRawParamTableForScopes{
				p: p,
			}
		}
	}
	return nil
}

func (script *addRawParamTableForScopes) Up(basicRes context.BasicRes) errors.Error {
	scope := script.p.(plugin.PluginSource).Scope()
	err := basicRes.GetDal().UpdateColumn(scope.TableName(), "_raw_data_table", fmt.Sprintf("_raw_%s_scopes", script.p.Name()),
		dal.Where("1=1"),
	)
	if err != nil {
		return err
	}
	return nil
}

func (*addRawParamTableForScopes) Version() uint64 {
	return 20230630000001
}

func (script *addRawParamTableForScopes) Name() string {
	return fmt.Sprintf("populated _raw_data_table column of plugin %s", script.p.Name())
}
