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

import "time"

type ApiBambooLink struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

type ApiBambooSizeData struct {
	Size       int `json:"size"`
	StartIndex int `json:"start-index"`
	MaxResult  int `json:"max-result"`
}
type ApiBambooKey struct {
	Key string `json:"key"`
}

type ApiBambooOperations struct {
	CanView                   bool `json:"canView"`
	CanEdit                   bool `json:"canEdit"`
	CanDelete                 bool `json:"canDelete"`
	AllowedToExecute          bool `json:"allowedToExecute"`
	CanExecute                bool `json:"canExecute"`
	AllowedToCreateVersion    bool `json:"allowedToCreateVersion"`
	AllowedToSetVersionStatus bool `json:"allowedToSetVersionStatus"`
}

func unixForBambooDeployBuild(time_unix int64) *time.Time {
	t := time.Unix(time_unix/1000, 0)
	return &t
}
