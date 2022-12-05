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

package helper

import (
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"regexp"
)

// RegexEnricher process value with regex pattern
type RegexEnricher struct {
	// This field will store compiled regular expression for every pattern
	regexpMap map[string]*regexp.Regexp
}

// NewRegexEnricher initialize a regexEnricher
func NewRegexEnricher() *RegexEnricher {
	return &RegexEnricher{regexpMap: make(map[string]*regexp.Regexp)}
}

// AddRegexp will add compiled regular expression for pattern to regexpMap
func (r *RegexEnricher) AddRegexp(patterns ...string) errors.Error {
	for _, pattern := range patterns {
		if len(pattern) > 0 {
			regex, err := errors.Convert01(regexp.Compile(pattern))
			if err != nil {
				return errors.Default.Wrap(err, fmt.Sprintf("Fail to compile pattern for regex pattern: %s", pattern))
			}
			r.regexpMap[pattern] = regex
		}
	}
	return nil
}

// GetEnrichResult will get compiled regular expression from map by pattern,
// and check if v matches compiled regular expression,
// lastly, will return corresponding value(result or empty)
func (r *RegexEnricher) GetEnrichResult(pattern string, v string, result string) string {
	if result == devops.PRODUCTION && pattern == "" {
		return result
	}
	regex := r.regexpMap[pattern]
	if regex != nil {
		if flag := regex.FindString(v); flag != "" {
			return result
		}
	}
	return ""
}
