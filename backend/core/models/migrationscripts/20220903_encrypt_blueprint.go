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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*encryptBlueprint)(nil)

type encryptBlueprint struct{}

type BlueprintEncryption0904 struct {
	archived.Model
	Plan     string
	Settings string
}

func (script *encryptBlueprint) Up(basicRes context.BasicRes) errors.Error {
	encKey := basicRes.GetConfig(plugin.EncodeKeyEnvStr)
	if encKey == "" {
		return errors.BadInput.New("invalid encKey")
	}
	err := migrationhelper.TransformColumns(
		basicRes,
		script,
		"_devlake_blueprints",
		[]string{"plan", "settings"},
		func(src *BlueprintEncryption0904) (*BlueprintEncryption0904, errors.Error) {
			plan, err := plugin.Encrypt(encKey, src.Plan)
			if err != nil {
				return nil, err
			}
			settings, err := plugin.Encrypt(encKey, src.Settings)
			if err != nil {
				return nil, err
			}
			return &BlueprintEncryption0904{
				Model:    src.Model,
				Plan:     plan,
				Settings: settings,
			}, nil
		},
	)
	return err
}

func (*encryptBlueprint) Version() uint64 {
	return 20220904142321
}

func (*encryptBlueprint) Name() string {
	return "encrypt Blueprint"
}
