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
	"github.com/apache/incubator-devlake/plugins/kube_deployment/models/migrationscripts/archived"
)

type addInitTables struct{}

func (u *addInitTables) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	err := db.DropTables(
		&archived.KubeDeploymentRevision{},
		&archived.KubeDeployment{},
		&archived.KubeConnection{},
	)

	if err != nil {
		return err
	}

	err = migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.KubeDeploymentRevision{},
		&archived.KubeConnection{},
	)

	if err != nil {
		return err
	}

	// encodeKey := basicRes.GetConfig(plugin.EncodeKeyEnvStr)
	// connection := &archived.KubeConnection{}
	// connection.Endpoint = "https://api.coincap.io/" // path -> v2/assets
	// connection.Token = "1234567890"
	// connection.Name = `kube_deployment`

	// if connection.Endpoint != `` && connection.Token != `` && encodeKey != `` {
	// 	connection.Token, err = plugin.Encrypt(encodeKey, connection.Token)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// update from .env and save to db

	// 	err = db.CreateIfNotExist(connection)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (*addInitTables) Version() uint64 {
	return 20230518000001
}

func (*addInitTables) Name() string {
	return "kube deployment revision init schemas"
}
