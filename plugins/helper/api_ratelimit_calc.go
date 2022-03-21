package helper

import (
	"net/http"
	"time"
)

// A helper to calculate api rate limit dynamically, assuming api returning remaining/resettime information
type ApiRateLimitCalculator struct {
	UserRateLimitPerHour    int
	GlobalRateLimitPerHour  int
	MaxRetry                int
	Method                  string
	ApiPath                 string
	DynamicRateLimitPerHour func(res *http.Response) (int, time.Duration, error)
}

func (c *ApiRateLimitCalculator) Calculate(apiClient *ApiClient) (int, time.Duration, error) {
	// user specified rate limit has the highest priority
	if c.UserRateLimitPerHour > 0 {
		return c.UserRateLimitPerHour, 1 * time.Hour, nil
	}
	// plugin dynamical rate limit is medium priority
	if c.DynamicRateLimitPerHour != nil {
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
			return c.DynamicRateLimitPerHour(res)
		}
	}

	// global default rate limit is the lowest
	return c.GlobalRateLimitPerHour, 1 * time.Hour, nil
}
