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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
)

// CallDB wraps DB calls with this signature, and handles the case if the struct is wrapped in a models.DynamicTabler.
func CallDB(f func(any, ...dal.Clause) errors.Error, x any, clauses ...dal.Clause) errors.Error {
	if dynamic, ok := x.(*models.DynamicTabler); ok {
		clauses = append(clauses, dal.From(dynamic.TableName()))
		x = dynamic.Unwrap()
	}
	return f(x, clauses...)
}
