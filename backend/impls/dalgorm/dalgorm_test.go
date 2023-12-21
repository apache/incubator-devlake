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

package dalgorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateQuery(t *testing.T) {
	for _, target := range []string{
		"begin",
		" begin",
		"begin ",
		" begin ",
		"begin;",
		"begin ;",
		"begin ; ;",
		"BEGIN ; ;",
		"start transaction",
		"start   transaction",
		"START   TRANSACTION",
		"start\t\n   transaction",
		" ;; start transaction",
	} {
		assert.EqualError(t, validateQuery(target), "illegal invocation, use the `Begin()` method instead", "failed text: `%s`", target)
	}
	for _, target := range []string{
		"select 1",
		"update a set b = c",
	} {
		assert.Nil(t, validateQuery(target), "failed text: `%s`", target)
	}
}
