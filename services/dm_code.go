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

package services

import "github.com/apache/incubator-devlake/models/domainlayer/code"

// GetRepos FIXME ...
func GetRepos() ([]*code.Repo, int64, error) {
	repos := make([]*code.Repo, 0)
	db := db.Model(repos).Order("id DESC")
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Find(&repos).Error
	if err != nil {
		return nil, count, err
	}
	return repos, count, nil
}
