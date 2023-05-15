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

package api

import (
	"fmt"
	"regexp"

	"github.com/apache/incubator-devlake/core/errors"
)

// RegexEnricher process value with regex pattern
// TODO: remove Enricher from naming since it is more like a util function
type RegexEnricher struct {
	// This field will store compiled regular expression for every pattern
	regexpMap map[string]*regexp.Regexp
}

// NewRegexEnricher initialize a regexEnricher
func NewRegexEnricher() *RegexEnricher {
	return &RegexEnricher{regexpMap: make(map[string]*regexp.Regexp)}
}

// AddRegexp will add compiled regular expression for pattern to regexpMap
// TODO: to be removed
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
// TODO: to be removed
func (r *RegexEnricher) GetEnrichResult(pattern string, v string, result string) string {
	if pattern == "" {
		return ""
	}
	regex := r.regexpMap[pattern]
	if regex != nil {
		if flag := regex.FindString(v); flag != "" {
			return result
		}
	}
	return ""
}

// TryAdd a named regexp if given pattern is not empty
func (r *RegexEnricher) TryAdd(name, pattern string) errors.Error {
	if pattern == "" {
		return nil
	}
	if _, ok := r.regexpMap[name]; ok {
		return errors.Default.New(fmt.Sprintf("Regex pattern with name: %s already exists", name))
	}
	regex, err := errors.Convert01(regexp.Compile(pattern))
	if err != nil {
		return errors.BadInput.Wrap(err, fmt.Sprintf("Fail to compile pattern for regex pattern: %s", pattern))
	}
	r.regexpMap[name] = regex
	return nil
}

// ReturnNameIfMatched will return name if any of the targets matches the regex associated with the given name
func (r *RegexEnricher) ReturnNameIfMatched(name string, targets ...string) string {
	if regex, ok := r.regexpMap[name]; !ok {
		return ""
	} else {
		for _, target := range targets {
			if regex.MatchString(target) {
				return name
			}
		}
	}
	return ""
}

// ReturnNameIfMatchedOrOmitted returns the given name if regex of the given name is omitted or fallback to ReturnNameIfMatched
func (r *RegexEnricher) ReturnNameIfOmittedOrMatched(name string, targets ...string) string {
	if _, ok := r.regexpMap[name]; !ok {
		return name
	}
	return r.ReturnNameIfMatched(name, targets...)
}
