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

package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildUserActivityAllowedEmailSetNormalizesEmails(t *testing.T) {
	allowedEmails := buildUserActivityAllowedEmailSet([]string{
		" Alice@example.com ",
		"BOB@example.com",
		"",
	})

	assert.Len(t, allowedEmails, 2)
	_, hasAlice := allowedEmails["alice@example.com"]
	_, hasBob := allowedEmails["bob@example.com"]
	assert.True(t, hasAlice)
	assert.True(t, hasBob)
}

func TestShouldExtractUserActivityEmail(t *testing.T) {
	allowedEmails := buildUserActivityAllowedEmailSet([]string{"alice@example.com"})

	assert.True(t, shouldExtractUserActivityEmail(normalizeUserActivityEmail(" Alice@example.com "), allowedEmails))
	assert.False(t, shouldExtractUserActivityEmail(normalizeUserActivityEmail("bob@example.com"), allowedEmails))
	assert.False(t, shouldExtractUserActivityEmail("", allowedEmails))
}
