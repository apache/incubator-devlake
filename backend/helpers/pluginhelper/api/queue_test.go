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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueueIterator(t *testing.T) {
	it := NewQueueIterator()
	it.Push("a")
	it.Push("b")
	require.True(t, it.HasNext())
	folderRaw, err := it.Fetch()
	require.NoError(t, err)
	data := folderRaw.(string)
	require.Equal(t, "a", data)
	require.True(t, it.HasNext())
	folderRaw, err = it.Fetch()
	require.NoError(t, err)
	data = folderRaw.(string)
	require.Equal(t, "b", data)
	require.True(t, it.HasNext())

	empty := false
	waited := false
	go func() {
		data, err := it.Fetch()
		require.Equal(t, err, nil)
		require.Equal(t, data, nil)
		require.False(t, it.HasNext())
		empty = true
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		waited = true
		it.Finish(2)
	}()

	for !empty {
		time.Sleep(time.Millisecond)
	}

	require.True(t, waited)
}
