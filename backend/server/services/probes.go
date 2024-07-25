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

package services

import (
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
)

// Ready returns the readiness status of the service
func Ready() (string, errors.Error) {
	var err errors.Error
	if serviceStatus != SERVICE_STATUS_READY {
		err = errors.Unavailable.New("service is not ready: " + serviceStatus)
	}
	return serviceStatus, err
}

// Health returns the health status of the service
func Health() (string, errors.Error) {
	// return true, nil unless we are 100% sure that the service is unhealthy
	if serviceStatus != SERVICE_STATUS_READY {
		return "maybe", nil
	}
	// cover the cases #5711, #6685 that we ran into in the pass
	// it is healthy if we could read one record from the pipelines table in 5 seconds
	result := make(chan errors.Error, 1)
	go func() {
		result <- db.All(&models.Pipeline{}, dal.Limit(1))
	}()
	select {
	case <-time.After(5 * time.Second):
		return "timeouted", errors.Default.New("timeout reading from pipelines")
	case err := <-result:
		if err != nil {
			return "bad", err
		}
		return "good", nil
	}
}
