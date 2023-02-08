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

package archived

type StFileMetrics struct {
	DomainEntity
	ComponentKey             string  `json:"component_key"`
	Project                  string  `gorm:"index;type:varchar(255)"` //domain project key
	FileName                 string  `json:"file_name"`
	FilePath                 string  `json:"file_path"`
	FileLanguage             string  `json:"file_language"`
	BatchID                  string  `json:"batch_id"`
	CodeSmells               int     `json:"code_smells"`
	SqaleIndex               string  `json:"sqale_index"`
	SqaleRating              string  `json:"sqale_rating"`
	Bugs                     int     `json:"bugs"`
	ReliabilityRating        string  `json:"reliability_rating"`
	Vulnerabilities          int     `json:"vulnerabilities"`
	SecurityRating           string  `json:"security_rating"`
	SecurityHotspots         int     `json:"security_hotspots"`
	SecurityHotspotsReviewed float64 `json:"security_hotspots_reviewed"`
	SecurityReviewRating     string  `json:"security_review_rating"`
	Ncloc                    int     `json:"ncloc"`
	Coverage                 float64 `json:"coverage"`
	LinesToCover             int     `json:"lines_to_cover"`
	DuplicatedLinesDensity   float64 `json:"duplicated_lines_density"`
	DuplicatedBlocks         int     `json:"duplicated_blocks"`
}

func (StFileMetrics) TableName() string {
	return "st_file_metrics"
}
