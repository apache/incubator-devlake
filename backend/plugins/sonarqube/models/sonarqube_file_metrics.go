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
	"github.com/apache/incubator-devlake/core/models/common"
)

type SonarqubeFileMetrics struct {
	ConnectionId             uint64  `gorm:"primaryKey"`
	FileMetricsKey           string  `json:"componentKey" gorm:"primaryKey"`
	ProjectKey               string  `json:"projectKey" gorm:"index"`
	BatchID                  string  `json:"batchId"`
	FileName                 string  `json:"fileName"`
	FilePath                 string  `json:"filePath"`
	FileLanguage             string  `json:"fileLanguage"`
	CodeSmells               int     `json:"codeSmells"`
	SqaleIndex               string  `json:"sqaleIndex"`
	SqaleRating              string  `json:"sqaleRating"`
	Bugs                     int     `json:"bugs"`
	ReliabilityRating        string  `json:"reliabilityRating"`
	Vulnerabilities          int     `json:"vulnerabilities"`
	SecurityRating           string  `json:"securityRating"`
	SecurityHotspots         int     `json:"securityHotspots"`
	SecurityHotspotsReviewed float64 `json:"securityHotspotsReviewed"`
	SecurityReviewRating     string  `json:"securityReviewRating"`
	Ncloc                    int     `json:"ncloc"`
	Coverage                 float64 `json:"coverage"`
	LinesToCover             int     `json:"linesToCover"`
	DuplicatedLinesDensity   float64 `json:"duplicatedLinesDensity"`
	DuplicatedBlocks         int     `json:"duplicatedBlocks"`
	common.NoPKModel
}

func (SonarqubeFileMetrics) TableName() string {
	return "_tool_sonarqube_file_metrics"
}
