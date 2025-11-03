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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectContainerImages_ReturnsSortedUniqueImages(t *testing.T) {
	op := &ArgocdApiSyncOperation{}
	op.Metadata.Images = []string{"registry.example.com/system:ops", " registry.example.com/sidecar:def456 "}
	op.Metadata.Resources = []ArgocdApiSyncOperationMetadataResource{
		{Images: []string{"registry.example.com/app:abc123", "", "registry.example.com/api:789xyz"}},
	}
	op.Operation.Metadata.Images = []string{"registry.example.com/worker:alpha"}
	op.Operation.Metadata.Resources = []ArgocdApiSyncOperationMetadataResource{
		{Images: []string{"registry.example.com/rollout:blue"}},
	}
	op.Operation.Sync.Resources = []ArgocdApiSyncResourceItem{
		{Images: []string{"registry.example.com/canary:latest"}},
	}
	op.SyncResult.Resources = []ArgocdApiSyncResourceItem{
		{Images: []string{"registry.example.com/worker:alpha", "registry.example.com/api:789xyz"}},
	}

	images := collectContainerImages(op)

	expected := []string{
		"registry.example.com/api:789xyz",
		"registry.example.com/app:abc123",
		"registry.example.com/canary:latest",
		"registry.example.com/rollout:blue",
		"registry.example.com/sidecar:def456",
		"registry.example.com/system:ops",
		"registry.example.com/worker:alpha",
	}
	assert.Equal(t, expected, images)
}

func TestCollectContainerImages_EmptyInputReturnsNil(t *testing.T) {
	op := &ArgocdApiSyncOperation{}

	assert.Nil(t, collectContainerImages(op))
	assert.Nil(t, collectContainerImages(nil))
}

func TestNormalizeImages(t *testing.T) {
	input := []string{" registry.example.com/app:1.0 ", "registry.example.com/app:1.0", "registry.example.com/api:2.0", ""}
	expected := []string{"registry.example.com/api:2.0", "registry.example.com/app:1.0"}
	assert.Equal(t, expected, normalizeImages(input))
}

// Fallback: no images in payload → expect none here (other fallbacks tested elsewhere)
func TestCollectContainerImages_FallbackRevisionAndSummary(t *testing.T) {
	revision := "abcdef1234567890"
	apiPayload := ArgocdApiSyncOperation{Revision: revision}
	assert.Nil(t, collectContainerImages(&apiPayload))

	// normalizeImages: dedupe + sort
	assert.Equal(t, []string{"a", "b"}, normalizeImages([]string{"b", "a", "b"}))
}
