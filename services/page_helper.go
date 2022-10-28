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

import "gorm.io/gorm"

const maxPageSize = 100

func processDbClausesWithPager(tx *gorm.DB, pageSize int, page int) *gorm.DB {
	if pageSize <= 0 || pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	tx = tx.Limit(pageSize)

	if page > 0 {
		offset := pageSize * (page - 1)
		tx = tx.Offset(offset)
	}
	return tx
}
