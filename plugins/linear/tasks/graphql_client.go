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

package tasks

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/linear/models"
	"github.com/merico-dev/graphql"
)

type GraphQueryRateLimit struct {
	RateLimit struct {
		Limit     graphql.Int
		Remaining graphql.Int
		ResetAt   time.Time
	}
}

func NewLinearGraphqlClient(taskCtx core.TaskContext, connection *models.LinearConnection) (*helper.GraphqlAsyncClient, errors.Error) {
	// create synchronize api client so we can calculate api rate limit dynamically
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", connection.Token),
	}
	apiClient, err := helper.NewApiClient(taskCtx.GetContext(), connection.Endpoint, headers, 0, connection.Proxy, taskCtx)
	if err != nil {
		return nil, err
	}
	apiClient.SetAfterFunction(func(res *http.Response) errors.Error {
		if res.StatusCode == http.StatusUnauthorized {
			return errors.Default.New(fmt.Sprintf("authentication failed, please check your AccessToken"))
		}
		return nil
	})

	graphqlClient := graphql.NewClient(connection.Endpoint+`graphql`, apiClient.GetHttpClient())

	asyncGraphqlClient := helper.CreateAsyncGraphqlClient(taskCtx.GetContext(), graphqlClient, taskCtx.GetLogger(),
		func(ctx context.Context, client *graphql.Client, logger core.Logger) (rateRemaining int, resetAt *time.Time, err errors.Error) {
			var query GraphQueryRateLimit
			err = errors.Convert(client.Query(taskCtx.GetContext(), &query, nil))
			if err != nil {
				return 0, nil, err
			}
			logger.Info(`linear graphql init success with remaining %d/%d and will reset at %s`,
				query.RateLimit.Remaining, query.RateLimit.Limit, query.RateLimit.ResetAt)
			return int(query.RateLimit.Remaining), &query.RateLimit.ResetAt, nil
		})

	asyncGraphqlClient.SetGetRateCost(func(q interface{}) int {
		v := reflect.ValueOf(q)
		return int(v.Elem().FieldByName(`RateLimit`).FieldByName(`Cost`).Int())
	})

	return asyncGraphqlClient, nil
}
