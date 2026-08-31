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

package api

import (
	"github.com/go-playground/validator/v10"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/kiro/models"
)

var (
	vld              *validator.Validate
	connectionHelper *api.ConnectionApiHelper
	basicRes         context.BasicRes
	// Scope config is NoScopeConfig: the report CSV's meaning is fixed by AWS
	// and uniform across an organization, so there is nothing per-scope to
	// configure.
	dsHelper *api.DsHelper[models.KiroConnection, models.KiroS3Slice, srvhelper.NoScopeConfig]
)

func Init(br context.BasicRes, p plugin.PluginMeta) {
	basicRes = br
	vld = validator.New()
	connectionHelper = api.NewConnectionHelper(
		basicRes,
		vld,
		p.Name(),
	)

	dsHelper = api.NewDataSourceHelper[
		models.KiroConnection, models.KiroS3Slice, srvhelper.NoScopeConfig,
	](
		basicRes,
		p.Name(),
		// Searchable scope fields.
		[]string{"accountId", "name"},
		func(c models.KiroConnection) models.KiroConnection { return c.Sanitize() },
		func(s models.KiroS3Slice) models.KiroS3Slice { return s.Sanitize() },
		nil,
	)

	// Scope browsing and search are implemented directly in remote_api.go rather
	// than through the shared DsRemoteApiScopeList/Search helpers. Those build an
	// HTTP client from the connection first, and that constructor runs a DNS
	// check on the endpoint - a bucket name is not a hostname, so it fails. The
	// helpers assume an HTTP data source; this one is S3.
}
