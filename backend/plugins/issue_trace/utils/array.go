package utils

import "strings"

func StringContains(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}

	return false
}

func ResolveMultiChangelogs(from, to string) (removedFrom []string, addedTo []string) {
	splitFromItems := strings.Split(from, ",")
	splitToItems := strings.Split(to, ",")
	fromItems := make([]string, 0)
	toItems := make([]string, 0)
	for _, v := range splitFromItems {
		if strings.TrimSpace(v) == "" {
			continue
		}
		fromItems = append(fromItems, strings.TrimSpace(v))
	}
	for _, v := range splitToItems {
		if strings.TrimSpace(v) == "" {
			continue
		}
		toItems = append(toItems, strings.TrimSpace(v))
	}
	for _, v := range fromItems {
		if StringContains(toItems, v) {
			continue
		}
		removedFrom = append(removedFrom, v)
	}
	for _, v := range toItems {
		if StringContains(fromItems, v) {
			continue
		}
		addedTo = append(addedTo, v)
	}
	return
}
