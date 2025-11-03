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

import "github.com/apache/incubator-devlake/core/models/common"

// ArgocdRevisionImage captures the container images observed for a given
// Argo CD application revision. It enables historical lookups so that
// previously processed sync operations retain the images that were active
// when they first ran, even after subsequent deployments update the
// application summary images.
type ArgocdRevisionImage struct {
	ConnectionId    uint64   `gorm:"primaryKey"`
	ApplicationName string   `gorm:"primaryKey;type:varchar(255)"`
	Revision        string   `gorm:"primaryKey;type:varchar(255)"`
	Images          []string `gorm:"type:json;serializer:json"`
	common.NoPKModel
}

func (ArgocdRevisionImage) TableName() string {
	return "_tool_argocd_revision_images"
}
