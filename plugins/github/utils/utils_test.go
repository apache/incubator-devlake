// https://golang.org/doc/tutorial/add-a-test

package utils

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
