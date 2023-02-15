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
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type collectorLatestState20230213 struct {
	TimeAfter *time.Time
}

func (collectorLatestState20230213) TableName() string {
	return "_devlake_collector_latest_state"
}

type addTimeAfterToCollectorMeta20230213 struct{}

func (script *addTimeAfterToCollectorMeta20230213) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &collectorLatestState20230213{})
}

func (*addTimeAfterToCollectorMeta20230213) Version() uint64 {
	return 20230213200040
}

func (*addTimeAfterToCollectorMeta20230213) Name() string {
	return "add time_after to _devlake_collector_latest_state"
}
