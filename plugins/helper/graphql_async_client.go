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
	"context"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/merico-dev/graphql"
	"sync"
	"time"
)

// GraphqlAsyncClient send graphql one by one
type GraphqlAsyncClient struct {
	ctx          context.Context
	cancel       context.CancelFunc
	client       *graphql.Client
	logger       core.Logger
	mu           sync.Mutex
	waitGroup    sync.WaitGroup
	workerErrors []error

	rateExhaustCond  *sync.Cond
	rateRemaining    int
	getRateRemaining func(context.Context, *graphql.Client, core.Logger) (rateRemaining int, resetAt *time.Time, err errors.Error)
	getRateCost      func(q interface{}) int
}

// CreateAsyncGraphqlClient creates a new GraphqlAsyncClient
func CreateAsyncGraphqlClient(
	ctx context.Context,
	graphqlClient *graphql.Client,
	logger core.Logger,
	getRateRemaining func(context.Context, *graphql.Client, core.Logger) (rateRemaining int, resetAt *time.Time, err errors.Error),
) *GraphqlAsyncClient {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	graphqlAsyncClient := &GraphqlAsyncClient{
		ctx:              ctxWithCancel,
		cancel:           cancel,
		client:           graphqlClient,
		logger:           logger,
		rateExhaustCond:  sync.NewCond(&sync.Mutex{}),
		rateRemaining:    0,
		getRateRemaining: getRateRemaining,
	}
	if getRateRemaining != nil {
		rateRemaining, resetAt, err := getRateRemaining(ctx, graphqlClient, logger)
		if err != nil {
			panic(err)
		}
		graphqlAsyncClient.updateRateRemaining(rateRemaining, resetAt)
	}
	return graphqlAsyncClient
}

// updateRateRemaining call getRateRemaining to update rateRemaining periodically
func (apiClient *GraphqlAsyncClient) updateRateRemaining(rateRemaining int, resetAt *time.Time) {
	apiClient.rateRemaining = rateRemaining
	if rateRemaining > 0 {
		apiClient.rateExhaustCond.Signal()
	}
	go func() {
		nextDuring := 3 * time.Minute
		if resetAt != nil && resetAt.After(time.Now()) {
			nextDuring = time.Until(*resetAt)
		}
		<-time.After(nextDuring)
		newRateRemaining, newResetAt, err := apiClient.getRateRemaining(apiClient.ctx, apiClient.client, apiClient.logger)
		if err != nil {
			panic(err)
		}
		apiClient.updateRateRemaining(newRateRemaining, newResetAt)
	}()
}

// SetGetRateCost to calculate how many rate cost
// if not set, all query just cost 1
func (apiClient *GraphqlAsyncClient) SetGetRateCost(getRateCost func(q interface{}) int) {
	apiClient.getRateCost = getRateCost
}

// Query send a graphql request when get lock
func (apiClient *GraphqlAsyncClient) Query(q interface{}, variables map[string]interface{}) errors.Error {
	apiClient.waitGroup.Add(1)
	defer apiClient.waitGroup.Done()
	apiClient.mu.Lock()
	defer apiClient.mu.Unlock()

	apiClient.rateExhaustCond.L.Lock()
	defer apiClient.rateExhaustCond.L.Unlock()
	for apiClient.rateRemaining <= 0 {
		apiClient.logger.Info(`rate limit remaining exhausted, waiting for next period.`)
		apiClient.rateExhaustCond.Wait()
	}
	select {
	case <-apiClient.ctx.Done():
		return nil
	default:
		err := apiClient.client.Query(apiClient.ctx, q, variables)
		if err != nil {
			return errors.Default.Wrap(err, "error making GraphQL call")
		}
		cost := 1
		if apiClient.getRateCost != nil {
			cost = apiClient.getRateCost(q)
		}
		apiClient.rateRemaining -= cost
		apiClient.logger.Debug(`query cost %d in %v`, cost, variables)
		return nil
	}
}

// NextTick to return the NextTick of scheduler
func (apiClient *GraphqlAsyncClient) NextTick(task func() errors.Error) {
	// to make sure task will be enqueued
	apiClient.waitGroup.Add(1)
	go func() {
		select {
		case <-apiClient.ctx.Done():
			return
		default:
			go func() {
				// if set waitGroup done here, a serial of goruntine will block until son goruntine finish.
				// But if done out of this go func, so task will run after waitGroup finish
				// I have no idea about this now...
				defer apiClient.waitGroup.Done()
				apiClient.checkError(task())
			}()
		}
	}()
}

// Wait blocks until all async requests were done
func (apiClient *GraphqlAsyncClient) Wait() errors.Error {
	apiClient.waitGroup.Wait()
	if len(apiClient.workerErrors) > 0 {
		return errors.Default.Combine(apiClient.workerErrors)
	}
	return nil
}

func (apiClient *GraphqlAsyncClient) checkError(err errors.Error) {
	if err == nil {
		return
	}
	apiClient.workerErrors = append(apiClient.workerErrors, err)
}

// HasError return if any error occurred
func (apiClient *GraphqlAsyncClient) HasError() bool {
	return len(apiClient.workerErrors) > 0
}
