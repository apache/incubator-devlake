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

type pagerdutyConnection20230123 struct {
	Endpoint         string `mapstructure:"endpoint" env:"PAGERDUTY_ENDPOINT" validate:"required"`
	Proxy            string `mapstructure:"proxy" env:"PAGERDUTY_PROXY"`
	RateLimitPerHour int    `comment:"api request rate limit per hour"`
}

func (pagerdutyConnection20230123) TableName() string {
	return "_tool_pagerduty_connections"
}

type addPagerdutyConnectionFields20230123 struct{}

func (script *addPagerdutyConnectionFields20230123) Name() string {
	return "add connection config fields"
}

func (script *addPagerdutyConnectionFields20230123) Up(basicRes context.BasicRes) errors.Error {
	err := basicRes.GetDal().AutoMigrate(&pagerdutyConnection20230123{})
	if err != nil {
		return err
	}
	return nil
}

func (*addPagerdutyConnectionFields20230123) Version() uint64 {
	return 20230123000000
}
