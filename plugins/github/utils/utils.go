package utils

import (
	"errors"
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

func ConvertRateLimitInfo(date string, resetTime string, remaining string) (RateLimitInfo, error) {
	var rateLimitInfo RateLimitInfo
	var err error
	if date != "" {
		rateLimitInfo.Date, err = http.ParseTime(date)
		if err != nil {
			return rateLimitInfo, err
		}
	} else {
		return rateLimitInfo, errors.New("rate limit date was an empty string")
	}
	if resetTime != "" {
		resetInt, err := strconv.ParseInt(resetTime, 10, 64)
		if err != nil {
			return rateLimitInfo, err
		}
		rateLimitInfo.ResetTime = time.Unix(resetInt, 0)
	} else {
		return rateLimitInfo, errors.New("rate limit reset time was an empty string")
	}
	if remaining != "" {
		rateLimitInfo.Remaining, err = strconv.Atoi(remaining)
		if err != nil {
			return rateLimitInfo, err
		}
	} else {
		return rateLimitInfo, errors.New("rate remaining was an empty string")
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
func ConvertStringToInt(input string) (int, error) {
	return strconv.Atoi(input)
}
func GetPagingFromLinkHeader(link string) (PagingInfo, error) {
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
				return result, errors.New("parsed string values aren't long enough")
			}
		}
		return result, nil
	} else {
		return result, errors.New("the link string provided is invalid. There is likely no next page of data to fetch")
	}
}

func GetIssueIdByIssueUrl(s string) (int, error) {
	regex := regexp.MustCompile(`.*/issues/(\d+)`)
	groups := regex.FindStringSubmatch(s)
	if len(groups) > 0 {
		return strconv.Atoi(groups[1])
	} else {
		return 0, errors.New("invalid issue url")
	}
}
