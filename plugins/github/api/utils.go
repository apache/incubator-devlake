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

package api

import (
	"net/url"
	"strconv"
)

const pageSize = 50

func getPageParam(q url.Values) (int, int) {
	var size, page int
	if ps := q["pageSize"]; len(ps) > 0 {
		size, _ = strconv.Atoi(ps[0])
	}
	if p := q["page"]; len(p) > 0 {
		page, _ = strconv.Atoi(p[0])
	}
	if size < 1 {
		size = pageSize
	}
	if page < 1 {
		page = 1
	}
	return size, page
}

func getLimitOffset(q url.Values) (int, int) {
	limit, page := getPageParam(q)
	offset := (page - 1) * limit
	return limit, offset
}
