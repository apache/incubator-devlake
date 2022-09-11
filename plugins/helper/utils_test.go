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

// https://golang.org/doc/tutorial/add-a-test

package helper

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

// TestParseLinkHeader calls utils.TestParseLinkHeader with a Link header string, checking
// for a valid return value.
func TestParseLinkHeader(t *testing.T) {
	fmt.Println("INFO >>> Handles good link string")
	var pagingExpected = PagingInfo{
		Next:  15,
		Last:  34,
		First: 1,
		Prev:  13,
	}
	linkHeaderFull := `<https://api.github.com/search/code?q=addClass+user%3Amozilla&page=15>; rel="next",
  <https://api.github.com/search/code?q=addClass+user%3Amozilla&page=34>; rel="last",
  <https://api.github.com/search/code?q=addClass+user%3Amozilla&page=1>; rel="first",
  <https://api.github.com/search/code?q=addClass+user%3Amozilla&page=13>; rel="prev"`
	result, err := GetPagingFromLinkHeader(linkHeaderFull)
	if err != nil {
		fmt.Println("ERROR: could not get paging from link header", err)

	}
	assert.Equal(t, result, pagingExpected)
}
func TestParseLinkHeaderEmptyString(t *testing.T) {
	fmt.Println("INFO >>> Handles empty link string")
	var pagingExpected = PagingInfo{
		Next:  1,
		Last:  1,
		First: 1,
		Prev:  1,
	}
	linkHeaderFull := ``
	paginationInfo, _ := GetPagingFromLinkHeader(linkHeaderFull)

	assert.Equal(t, paginationInfo, pagingExpected)
}

// This test is incomplete.
func TestGetRateLimitPerSecond(t *testing.T) {
	date := "Mon, 20 Sep 2021 18:08:38 GMT"
	resetTime := "1632164442"
	remaining := "100000"

	rateLimitInfo, err := ConvertRateLimitInfo(date, resetTime, remaining)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	rateLimitPerSecond := GetRateLimitPerSecond(rateLimitInfo)
	assert.Equal(t, rateLimitPerSecond, 31)
}

func TestGetIssueIdByIssueUrl(t *testing.T) {
	s := "https://api.github.com/repos/octocat/Hello-World/issues/1347"
	s1, err := GetIssueIdByIssueUrl(s)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Equal(t, s1, 1347)
}
