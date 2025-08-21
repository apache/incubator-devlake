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

package utils

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

// RandLetterBytes returns a cryptographically secure random string with given length n
func RandLetterBytes(n int) (string, errors.Error) {
	if n < 0 {
		return "", errors.Default.New("n must be greater than 0")
	}
	ret := make([]byte, n)
	bi := big.NewInt(int64(len(letterBytes)))
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, bi)
		if err != nil {
			return "", errors.Convert(err)
		}
		ret[i] = letterBytes[num.Int64()]
	}

	return string(ret), nil
}

func SanitizeString(s string) string {
	if s == "" {
		return s
	}
	strLen := len(s)
	if strLen <= 2 {
		return strings.Repeat("*", strLen)
	}
	prefixLen, suffixLen := 2, 2
	if strLen <= 5 {
		prefixLen, suffixLen = 1, 1
	}
	return strings.Replace(s, s[prefixLen:strLen-suffixLen], strings.Repeat("*", strLen-prefixLen-suffixLen), -1)
}

// from https://stackoverflow.com/questions/12311033/extracting-substrings-in-go
func Substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}
