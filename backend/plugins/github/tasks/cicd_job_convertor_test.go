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
	"fmt"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func GenJobIDWithReflect(jobIdGen *didgen.DomainIdGenerator) {
	connectionId := uint64(1)
	runId := 1
	lineId := 1
	jobIdGen.Generate(connectionId, runId, lineId)
}

func GenJobID() {
	connectionId := uint64(1)
	runId := 1
	lineId := 1
	fmt.Sprintf("GithubJob:%d:%d:%d", connectionId, runId, lineId)
}

func BenchmarkGenJobIDWithReflect(b *testing.B) {
	mockMeta := mockplugin.NewPluginMeta(b)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/github")
	mockMeta.On("Name").Return("github").Maybe()
	err := plugin.RegisterPlugin("github", mockMeta)
	assert.NoError(b, err)

	jobIdGen := didgen.NewDomainIdGenerator(&models.GithubJob{})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GenJobIDWithReflect(jobIdGen)
	}
	b.StopTimer()
	//BenchmarkGenJobIDWithReflect-8   	 5611773	       208.9 ns/op
}

func BenchmarkGenJobID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenJobID()
	}
	//BenchmarkGenJobID-8   	11078593	        99.43 ns/op
}
