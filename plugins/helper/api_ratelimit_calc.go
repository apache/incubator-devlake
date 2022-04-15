package helper

import (
	"net/http"
	"time"
)

// A helper to calculate api rate limit dynamically, assuming api returning remaining/resettime information
type ApiRateLimitCalculator struct {
	UserRateLimitPerHour   int
	GlobalRateLimitPerHour int
	MaxRetry               int
	Method                 string
	ApiPath                string
	DynamicRateLimit       func(res *http.Response) (int, time.Duration, error)
}

func (c *ApiRateLimitCalculator) Calculate(apiClient *ApiClient) (int, time.Duration, error) {
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

		var err error
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
