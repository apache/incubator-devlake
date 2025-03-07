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
)

type expandTimelineRecordName struct{}

func (*expandTimelineRecordName) Up(baseRes context.BasicRes) errors.Error {
	return baseRes.GetDal().ModifyColumnType("_tool_azuredevops_go_timeline_records", "name", "text")
}

func (*expandTimelineRecordName) Version() uint64 {
	return 20250307100000
}

func (*expandTimelineRecordName) Name() string {
	return "make timeline record name column larger to support longer names"
}
