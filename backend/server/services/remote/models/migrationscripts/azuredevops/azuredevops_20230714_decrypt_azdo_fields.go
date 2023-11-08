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

package azuredevops

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*DecryptConnectionFields)(nil)

type azureDevopsConnection20230825 struct {
	archived.Model
	Name         string
	Token        string
	Proxy        *string
	Organization *string
}

type DecryptConnectionFields struct{}

func (script *DecryptConnectionFields) Up(basicRes context.BasicRes) errors.Error {
	encryptionSecret := basicRes.GetConfig(plugin.EncodeKeyEnvStr)
	if encryptionSecret == "" {
		return errors.BadInput.New("invalid encryptionSecret")
	}

	err := migrationhelper.TransformColumns(
		basicRes,
		script,
		"_tool_azuredevops_azuredevopsconnections",
		[]string{"name", "proxy", "organization"},
		func(src *azureDevopsConnection20230825) (*azureDevopsConnection20230825, errors.Error) {
			encName := src.Name
			name, err := plugin.Decrypt(encryptionSecret, encName)
			if err != nil {
				return src, nil
			}
			src.Name = name

			if src.Proxy != nil {
				encProxy := *src.Proxy
				decProxy, err := plugin.Decrypt(encryptionSecret, encProxy)
				if err != nil {
					return src, nil
				}
				if decProxy == "" {
					src.Proxy = nil
				} else {
					src.Proxy = &decProxy
				}
			}

			if src.Organization != nil {
				encOrg := *src.Organization
				decOrg, err := plugin.Decrypt(encryptionSecret, encOrg)
				if err != nil {
					return src, nil
				}
				if decOrg == "" {
					src.Organization = nil
				} else {
					src.Organization = &decOrg
				}
			}
			return src, nil
		},
	)
	return err
}

func (*DecryptConnectionFields) Version() uint64 {
	return 20230825090504
}

func (script *DecryptConnectionFields) Name() string {
	return "Decrypt Azure DevOps connection fields"
}
