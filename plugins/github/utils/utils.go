package utils

import (
	"errors"
	"regexp"
	"strings"
)

type PagingInfo struct {
	Next  string
	Last  string
	First string
	Prev  string
}

func GetPagingFromLinkHeader(link string) (PagingInfo, error) {
	result := PagingInfo{
		Next:  "",
		Last:  "",
		Prev:  "",
		First: "",
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
				switch pageNameString {
				case "next":
					result.Next = pageNumberString

				case "first":
					result.First = pageNumberString

				case "last":
					result.Last = pageNumberString

				case "prev":
					result.Prev = pageNumberString
				}

			} else {
				return result, errors.New("Parsed string values aren't long enough.")
			}
		}
		return result, nil
	} else {
		return result, errors.New("Length of the split array is too short.")
	}
}
