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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

// addRawDataOriginToCopilotSeats ensures _tool_copilot_seats includes RawDataOrigin columns.
// This is required by ApiExtractor/StatefulApiExtractor, which attach provenance to extracted records.
type addRawDataOriginToCopilotSeats struct{}

func (script *addRawDataOriginToCopilotSeats) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&models.GhCopilotSeat{},
	)
}

func (*addRawDataOriginToCopilotSeats) Version() uint64 {
	return 20260105000000
}

func (*addRawDataOriginToCopilotSeats) Name() string {
	return "copilot add raw data origin to seats"
}
