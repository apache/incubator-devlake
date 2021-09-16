package utils

import (
	"regexp"
	"strings"
)

type Paging struct {
	next  string
	last  string
	first string
	prev  string
}

func GetPagingFromLinkHeader(link string) (Paging, error) {
	var result Paging
	linksArray := strings.Split(link, ",")
	pattern1 := regexp.MustCompile(`page=*[0-9]+`)
	pattern2 := regexp.MustCompile(`rel="*[a-z]+`)
	for i := 0; i < len(linksArray); i++ {
		content := []byte(linksArray[i])
		loc1 := pattern1.FindIndex(content)
		loc2 := pattern2.FindIndex(content)
		pageNumberSubstring := string(content[loc1[0]:loc1[1]])
		pageNameSubstring := string(content[loc2[0]:loc2[1]])
		pageNumberString := strings.Replace(pageNumberSubstring, `page=`, ``, 1)
		pageNameString := strings.Replace(pageNameSubstring, `rel="`, ``, 1)
		switch pageNameString {
		case "next":
			result.next = pageNumberString

		case "first":
			result.first = pageNumberString

		case "last":
			result.last = pageNumberString

		case "prev":
			result.prev = pageNumberString
		}
	}
	return result, nil
}
