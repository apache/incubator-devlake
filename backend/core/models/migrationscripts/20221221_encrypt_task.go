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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*encryptTask221221)(nil)

type encryptTask221221 struct{}

type srcTaskEncryption221221 struct {
	archived.Model
	Options json.RawMessage
}

type dstTaskEncryption221221 struct {
	archived.Model
	Options string
}

func (script *encryptTask221221) Up(basicRes context.BasicRes) errors.Error {
	encKey := basicRes.GetConfig(plugin.EncodeKeyEnvStr)
	if encKey == "" {
		return errors.BadInput.New("invalid encKey")
	}
	err := migrationhelper.TransformColumns(
		basicRes,
		script,
		"_devlake_tasks",
		[]string{"options"},
		func(src *srcTaskEncryption221221) (*dstTaskEncryption221221, errors.Error) {
			options, err := plugin.Encrypt(encKey, string(src.Options))
			if err != nil {
				return nil, err
			}
			return &dstTaskEncryption221221{
				Model:   src.Model,
				Options: options,
			}, nil
		},
	)
	return err
}

func (*encryptTask221221) Version() uint64 {
	return 20221221162121
}

func (*encryptTask221221) Name() string {
	return "encrypt task.options"
}
