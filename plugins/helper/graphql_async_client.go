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
	"fmt"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/shurcooL/graphql"
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
	getRateRemaining func(context.Context, *graphql.Client, core.Logger) (rateRemaining int, resetAt *time.Time)
	getRateCost      func(q interface{}) int
}

// CreateAsyncGraphqlClient creates a new GraphqlAsyncClient
func CreateAsyncGraphqlClient(
	ctx context.Context,
	graphqlClient *graphql.Client,
	logger core.Logger,
	getRateRemaining func(context.Context, *graphql.Client, core.Logger) (rateRemaining int, resetAt *time.Time),
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
		rateRemaining, resetAt := getRateRemaining(ctx, graphqlClient, logger)
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
			nextDuring = resetAt.Sub(time.Now())
		}
		<-time.After(nextDuring)
		newRateRemaining, newResetAt := apiClient.getRateRemaining(apiClient.ctx, apiClient.client, apiClient.logger)
		apiClient.updateRateRemaining(newRateRemaining, newResetAt)
	}()
}

// SetGetRateCost to calculate how many rate cost
// if not set, all query just cost 1
func (apiClient *GraphqlAsyncClient) SetGetRateCost(getRateCost func(q interface{}) int) {
	apiClient.getRateCost = getRateCost
}

// Query send a graphql request when get lock
func (apiClient *GraphqlAsyncClient) Query(q interface{}, variables map[string]interface{}) error {
	apiClient.mu.Lock()
	defer apiClient.mu.Unlock()

	apiClient.rateExhaustCond.L.Lock()
	defer apiClient.rateExhaustCond.L.Unlock()
	for apiClient.rateRemaining <= 0 {
		apiClient.rateExhaustCond.Wait()
	}
	select {
	case <-apiClient.ctx.Done():
		return nil
	default:
		err := apiClient.client.Query(apiClient.ctx, q, variables)
		if err != nil {
			return err
		}
		cost := 1
		if apiClient.getRateCost != nil {
			cost = apiClient.getRateCost(q)
		}
		apiClient.rateRemaining -= cost
		return nil
	}
}

// NextTick to return the NextTick of scheduler
func (apiClient *GraphqlAsyncClient) NextTick(task func() error) {
	// to make sure task will be enqueued
	apiClient.waitGroup.Add(1)
	go func() {
		defer apiClient.waitGroup.Done()
		select {
		case <-apiClient.ctx.Done():
			return
		default:
			apiClient.checkError(task())
		}
	}()
}

// WaitAsync blocks until all async requests were done
func (apiClient *GraphqlAsyncClient) Wait() error {
	apiClient.waitGroup.Wait()
	if len(apiClient.workerErrors) > 0 {
		return fmt.Errorf("%s", apiClient.workerErrors)
	}
	return nil
}

func (apiClient *GraphqlAsyncClient) checkError(err error) {
	if err == nil {
		return
	}
	apiClient.workerErrors = append(apiClient.workerErrors, err)
}
