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

import "github.com/apache/incubator-devlake/core/errors"

// Pagination holds the paginate information
type Pagination struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

// GetPage returns current page number
func (p *Pagination) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

// GetPageSize returns a sensible page size based on input
func (p *Pagination) GetPageSize() int {
	return p.GetPageSizeOr(50)
}

// GetPageSizeOr returns the page size or fallback to `defaultVal`
func (p *Pagination) GetPageSizeOr(defaultVal int) int {
	if p.PageSize < 1 {
		return defaultVal
	}
	return p.PageSize
}

// GetSkip returns how many records  should be skipped for specified page
func (p *Pagination) GetSkip() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// VerifyStruct verifies given struct with `validator`
func VerifyStruct(v interface{}) errors.Error {
	err := vld.Struct(v)
	if err != nil {
		return errors.BadInput.Wrap(err, "data verification failed")
	}
	return nil
}
