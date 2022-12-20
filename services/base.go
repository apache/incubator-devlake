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

import "github.com/apache/incubator-devlake/errors"

type Pagination struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

func (p *Pagination) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

func (p *Pagination) GetPageSize() int {
	return p.GetPageSizeOr(50)
}

func (p *Pagination) GetPageSizeOr(defaultVal int) int {
	if p.PageSize < 1 {
		return defaultVal
	}
	return p.PageSize
}

func (p *Pagination) GetSkip() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

func VerifyStruct(v interface{}) errors.Error {
	err := vld.Struct(v)
	if err != nil {
		return errors.BadInput.Wrap(err, "data verification failed")
	}
	return nil
}
