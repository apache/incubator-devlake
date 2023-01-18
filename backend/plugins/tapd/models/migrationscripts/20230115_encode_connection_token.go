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
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/tapd/models/migrationscripts/archived"
)

type encodeConnToken struct{}

func (script *encodeConnToken) Up(basicRes context.BasicRes) errors.Error {
	encKey := config.GetConfig().GetString(plugin.EncodeKeyEnvStr)
	if encKey == "" {
		return errors.BadInput.New("jira v0.11 invalid encKey")
	}

	return migrationhelper.TransformColumns(
		basicRes,
		script,
		"_tool_tapd_connections",
		[]string{},
		func(old *archived.TapdConnection) (*archived.TapdConnection, errors.Error) {
			var err errors.Error
			conn := &archived.TapdConnection{}
			conn.ID = old.ID
			conn.Name = old.Name
			conn.Endpoint = old.Endpoint
			conn.Proxy = old.Proxy
			conn.RateLimitPerHour = old.RateLimitPerHour
			conn.Username = old.Username
			conn.Password, err = plugin.Encrypt(encKey, old.Password)
			if err != nil {
				return nil, err
			}
			return conn, nil
		})
}

func (*encodeConnToken) Version() uint64 {
	return 20230115201138
}

func (*encodeConnToken) Name() string {
	return "Tapd encode connection token"
}
