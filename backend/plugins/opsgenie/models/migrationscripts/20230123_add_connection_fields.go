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

type opsgenieConnection20230905 struct {
	Endpoint         string `mapstructure:"endpoint" env:"OPSGENIE_ENDPOINT" validate:"required"`
	Proxy            string `mapstructure:"proxy" env:"OPSGENIE_PROXY"`
	RateLimitPerHour int    `comment:"api request rate limit per hour"`
}

func (opsgenieConnection20230905) TableName() string {
	return "_tool_opsgenie_connections"
}

type addOpsgenieConnectionFields20230905 struct{}

func (script *addOpsgenieConnectionFields20230905) Name() string {
	return "add connection config fields"
}

func (script *addOpsgenieConnectionFields20230905) Up(basicRes context.BasicRes) errors.Error {
	err := basicRes.GetDal().AutoMigrate(&opsgenieConnection20230905{})
	if err != nil {
		return err
	}
	return nil
}

func (*addOpsgenieConnectionFields20230905) Version() uint64 {
	return 20230905000000
}
