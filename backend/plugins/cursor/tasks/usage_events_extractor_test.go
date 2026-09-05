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
	"encoding/json"
	"testing"
)

func TestUsageEventRecordUnmarshalFractionalRequestsCosts(t *testing.T) {
	raw := []byte(`{
		"timestamp":"1783617259310",
		"model":"composer-2.5-fast",
		"kind":"Usage-based",
		"requestsCosts":3.8,
		"userEmail":"user@example.com",
		"chargedCents":8,
		"conversationId":"846952e0-0221-473b-bb1d-e5d016f071b1"
	}`)

	var record usageEventRecord
	if err := json.Unmarshal(raw, &record); err != nil {
		t.Fatalf("unmarshal usage event: %v", err)
	}
	if record.RequestsCosts != 3.8 {
		t.Fatalf("requestsCosts = %v, want 3.8", record.RequestsCosts)
	}
}

func TestComputeEventIdWithFractionalRequestsCosts(t *testing.T) {
	fractional := computeEventId(
		"1783617259310",
		"user@example.com",
		"846952e0-0221-473b-bb1d-e5d016f071b1",
		"composer-2.5-fast",
		8,
		3.8,
	)
	integer := computeEventId(
		"1783617259310",
		"user@example.com",
		"846952e0-0221-473b-bb1d-e5d016f071b1",
		"composer-2.5-fast",
		8,
		2,
	)
	if fractional == integer {
		t.Fatal("expected fractional requestsCosts to change event id")
	}
}
