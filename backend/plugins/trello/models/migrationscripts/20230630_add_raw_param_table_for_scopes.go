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
	"encoding/json"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/trello/models/migrationscripts/archived"
)

var _ plugin.MigrationScript = (*addRawParamTableForScope)(nil)

type scope20230630 struct {
	ConnectionId  uint64
	BoardId       string
	RawDataTable  string `gorm:"column:_raw_data_table"`
	RawDataParams string `gorm:"column:_raw_data_params"`
}

type params20230630 struct {
	ConnectionId uint64
	BoardId      string
}

type addRawParamTableForScope struct{}

func (script *addRawParamTableForScope) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.CopyTableColumns(basicRes,
		archived.TrelloBoard{}.TableName(),
		archived.TrelloBoard{}.TableName(),
		func(src *scope20230630) (*scope20230630, errors.Error) {
			src.RawDataTable = "_raw_trello_scopes"
			src.RawDataParams = string(errors.Must1(json.Marshal(&params20230630{
				ConnectionId: src.ConnectionId,
				BoardId:      src.BoardId,
			})))
			return src, nil
		})
}

func (*addRawParamTableForScope) Version() uint64 {
	return 20230630000001
}

func (script *addRawParamTableForScope) Name() string {
	return "populated _raw_data columns for trello boards"
}
