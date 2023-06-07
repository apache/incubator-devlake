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

package models

import (
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

// type AccessToken struct {
// 	ClientId     string `mapstructure:"clientId" validate:"required" json:"clientId" gorm:"serializer:encdec"`
// 	ClientSecret string `mapstructure:"ClientSecret" validate:"required" json:"ClientSecret" gorm:"serializer:encdec"`
// }

// KubeDeploymentConn holds the essential information to connect to the Crypto Asset API
type KubeConn struct {
	helper.RestConnection `mapstructure:",squash"`
	// helper.AccessToken    `mapstructure:",squash"`
	Credentials string `mapstructure:"credentials" validate:"required" json:"credentials" gorm:"serializer:encdec"`
}

// KubeDeploymentConnection holds KubeDeploymentConn plus ID/Name for database storage
type KubeConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	KubeConn              `mapstructure:",squash"`
}

func (KubeConnection) TableName() string {
	return "_tool_kube_connections"
}
