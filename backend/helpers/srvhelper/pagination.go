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

package srvhelper

type Pagination struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	SearchTerm string `json:"searchTerm"`
}

func (self *Pagination) GetPage() int {
	if self.Page > 0 {
		return self.Page
	}
	return 1
}

func (self *Pagination) GetLimit() int {
	if self.PageSize > 0 {
		return self.PageSize
	}
	return 100
}

func (self *Pagination) GetOffset() int {
	return (self.GetPage() - 1) * self.GetLimit()
}
