package utils

// StringsUniq returns a new String Slice contains deduped elements from `source`
func StringsUniq(source []string) []string {
	book := make(map[string]bool, len(source))
	target := make([]string, 0, len(source))
	for _, str := range source {
		if !book[str] {
			book[str] = true
			target = append(target, str)
		}
	}
	return target
}

// StringsContains checks if  `source` String Slice contains `target` string
func StringsContains(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}
