package utils

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/merico-dev/lake/logger"
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
	convertedDate, err := http.ParseTime(date)
	if err != nil {
		fmt.Println("Convert error: ", err)
	}
	resetInt, err := strconv.ParseInt(resetTime, 10, 64)
	if err != nil {
		fmt.Println("Convert error: ", err)
	}
	convertedResetTime := time.Unix(resetInt, 0)
	convertedRemaining, err := strconv.Atoi(remaining)
	if err != nil {
		logger.Error("Convert error: ", err)
	}
	return RateLimitInfo{
		Date:      convertedDate,
		ResetTime: convertedResetTime,
		Remaining: convertedRemaining,
	}, nil
}

func GetRateLimitPerSecond(info RateLimitInfo) int {
	return info.Remaining / int(info.ResetTime.Unix()-info.Date.Unix()) * 9 / 10
}
func ConvertStringToInt(input string) (int, error) {
	return strconv.Atoi(input)
}
func GetPagingFromLinkHeader(link string) (PagingInfo, error) {
	result := PagingInfo{
		Next:  0,
		Last:  0,
		Prev:  0,
		First: 0,
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
				return result, errors.New("Parsed string values aren't long enough.")
			}
		}
		return result, nil
	} else {
		return result, errors.New("The link string provided is invalid.")
	}
}
