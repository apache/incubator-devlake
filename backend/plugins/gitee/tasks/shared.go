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
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PagingInfo struct {
	Next  int
	Last  int
	First int
	Prev  int
}

type RateLimitInfo struct {
	Date      time.Time
	ResetTime time.Time
	Remaining int
}

type GiteeInput struct {
	Repo  string
	Owner string
	Iid   int
}

func GetTotalPagesFromResponse(res *http.Response, args *api.ApiCollectorArgs) (int, errors.Error) {
	total := res.Header.Get("X-Total-Pages")
	if total == "" {
		return 0, nil
	}
	totalInt, err := strconv.Atoi(total)
	if err != nil {
		return 0, errors.Convert(err)
	}
	return totalInt, nil
}

func GetRawMessageFromResponse(res *http.Response) ([]json.RawMessage, errors.Error) {
	var rawMessages []json.RawMessage

	if res == nil {
		return nil, errors.Default.New("res is nil")
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error reading response from %s", res.Request.URL.String()))
	}

	err = errors.Convert(json.Unmarshal(resBody, &rawMessages))
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error decoding response from %s: raw response: %s", res.Request.URL.String(), string(resBody)))
	}

	return rawMessages, nil
}

func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, Table string) (*api.RawDataSubTaskArgs, *GiteeTaskData) {
	data := taskCtx.GetData().(*GiteeTaskData)
	RawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Options: data.Options,
		Table:   Table,
	}
	return RawDataSubTaskArgs, data
}

func ConvertRateLimitInfo(date string, resetTime string, remaining string) (RateLimitInfo, errors.Error) {
	var rateLimitInfo RateLimitInfo
	var err errors.Error
	if date != "" {
		rateLimitInfo.Date, err = errors.Convert01(http.ParseTime(date))
		if err != nil {
			return rateLimitInfo, err
		}
	} else {
		return rateLimitInfo, errors.Default.New("rate limit date was an empty string")
	}
	if resetTime != "" {
		resetInt, err := errors.Convert01(strconv.ParseInt(resetTime, 10, 64))
		if err != nil {
			return rateLimitInfo, err
		}
		rateLimitInfo.ResetTime = time.Unix(resetInt, 0)
	} else {
		return rateLimitInfo, errors.Default.New("rate limit reset time was an empty string")
	}
	if remaining != "" {
		rateLimitInfo.Remaining, err = ConvertStringToInt(remaining)
		if err != nil {
			return rateLimitInfo, err
		}
	} else {
		return rateLimitInfo, errors.Default.New("rate remaining was an empty string")
	}
	return rateLimitInfo, nil
}

func GetRateLimitPerSecond(info RateLimitInfo) int {
	unixResetTime := info.ResetTime.Unix()
	unixNow := info.Date.Unix()
	timeBetweenNowAndReset := unixResetTime - unixNow
	// Adjust the remaining to be less then actual to avoid hitting the limit exactly.
	multiplier := 0.98
	adjustedRemaining := float64(info.Remaining) * multiplier
	return int(adjustedRemaining / float64(timeBetweenNowAndReset)) //* multiplier
}
func ConvertStringToInt(input string) (int, errors.Error) {
	return errors.Convert01(strconv.Atoi(input))
}
func GetPagingFromLinkHeader(link string) (PagingInfo, errors.Error) {
	result := PagingInfo{
		Next:  1,
		Last:  1,
		Prev:  1,
		First: 1,
	}
	linksArray := strings.Split(link, ",")
	pattern1 := regexp.MustCompile(`page=*[0-9]+`)
	pattern2 := regexp.MustCompile(`rel="*[a-z]+`)
	if len(linksArray) >= 2 {
		for i := 0; i < len(linksArray); i++ {
			content := []byte(linksArray[i])
			loc1 := pattern1.FindIndex(content)
			loc2 := pattern2.FindIndex(content)
			if len(loc1) >= 2 && len(loc2) >= 2 {
				pageNumberSubstring := string(content[loc1[0]:loc1[1]])
				pageNumberString := strings.Replace(pageNumberSubstring, `page=`, ``, 1)
				pageNameSubstring := string(content[loc2[0]:loc2[1]])
				pageNameString := strings.Replace(pageNameSubstring, `rel="`, ``, 1)

				pageNumberInt, convertErr := ConvertStringToInt(pageNumberString)
				if convertErr != nil {
					return result, convertErr
				}
				switch pageNameString {
				case "next":
					result.Next = pageNumberInt

				case "first":
					result.First = pageNumberInt

				case "last":
					result.Last = pageNumberInt

				case "prev":
					result.Prev = pageNumberInt
				}
			} else {
				return result, errors.Default.New("parsed string values aren't long enough")
			}
		}
		return result, nil
	} else {
		return result, errors.Default.New("the link string provided is invalid. There is likely no next page of data to fetch")
	}
}
