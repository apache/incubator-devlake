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
	"github.com/apache/incubator-devlake/errors"
	"net/http"
	"time"
)

// ApiRateLimitCalculator is A helper to calculate api rate limit dynamically, assuming api returning remaining/resettime information
type ApiRateLimitCalculator struct {
	UserRateLimitPerHour   int
	GlobalRateLimitPerHour int
	MaxRetry               int
	Method                 string
	ApiPath                string
	DynamicRateLimit       func(res *http.Response) (int, time.Duration, errors.Error)
}

// Calculate FIXME ...
func (c *ApiRateLimitCalculator) Calculate(apiClient *ApiClient) (int, time.Duration, errors.Error) {
	// user specified rate limit has the highest priority
	if c.UserRateLimitPerHour > 0 {
		return c.UserRateLimitPerHour, 1 * time.Hour, nil
	}
	// plugin dynamical rate limit is medium priority
	if c.DynamicRateLimit != nil {
		method := c.Method
		if method == "" {
			method = http.MethodOptions
		}

		var err errors.Error
		var res *http.Response
		for i := 0; i < c.MaxRetry; i++ {
			res, err = apiClient.Do(method, c.ApiPath, nil, nil, nil)
			if err != nil {
				continue
			}
			requests, duration, err := c.DynamicRateLimit(res)
			if err != nil {
				return requests, duration, err
			}

			if duration == 0 {
				return c.GlobalRateLimitPerHour, 1 * time.Hour, nil
			}
			requests = int(float32(requests) * 0.95)
			return requests, duration, err
		}
	}

	// global default rate limit is the lowest
	return c.GlobalRateLimitPerHour, 1 * time.Hour, nil
}
