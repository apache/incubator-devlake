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

package models

import (
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFolderInput(t *testing.T) {
	it := helper.NewQueueIterator()
	it.Push(NewFolderInput("a"))
	it.Push(NewFolderInput("b"))
	require.True(t, it.HasNext())
	folderRaw, err := it.Fetch()
	require.NoError(t, err)
	folder := folderRaw.(*FolderInput)
	require.Equal(t, "a", folder.Data())
	require.True(t, it.HasNext())
	folderRaw, err = it.Fetch()
	require.NoError(t, err)
	folder = folderRaw.(*FolderInput)
	require.Equal(t, "b", folder.Data())
	require.False(t, it.HasNext())
}
