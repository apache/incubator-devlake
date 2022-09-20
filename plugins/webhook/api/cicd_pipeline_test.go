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
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTypeAndResultFromTasks(t *testing.T) {
	testDomainTask := devops.CICDTask{Type: devops.TEST}
	buildDomainTask := devops.CICDTask{Type: devops.BUILD}
	deploymentDomainTask := devops.CICDTask{Type: devops.DEPLOYMENT}

	abortDomainTask := devops.CICDTask{Result: `ABORT`}
	failureDomainTask := devops.CICDTask{Result: `FAILURE`}
	successDomainTask := devops.CICDTask{Result: `SUCCESS`}

	pipelineType, _ := getTypeAndResultFromTasks([]devops.CICDTask{testDomainTask})
	assert.Equal(t, devops.TEST, pipelineType)

	pipelineType, _ = getTypeAndResultFromTasks([]devops.CICDTask{buildDomainTask})
	assert.Equal(t, devops.BUILD, pipelineType)

	pipelineType, _ = getTypeAndResultFromTasks([]devops.CICDTask{deploymentDomainTask, deploymentDomainTask, deploymentDomainTask})
	assert.Equal(t, devops.DEPLOYMENT, pipelineType)

	pipelineType, _ = getTypeAndResultFromTasks([]devops.CICDTask{buildDomainTask, deploymentDomainTask, buildDomainTask})
	assert.Equal(t, ``, pipelineType)

	_, result := getTypeAndResultFromTasks([]devops.CICDTask{abortDomainTask, failureDomainTask, successDomainTask})
	assert.Equal(t, `ABORT`, result)

	_, result = getTypeAndResultFromTasks([]devops.CICDTask{failureDomainTask, successDomainTask})
	assert.Equal(t, `FAILURE`, result)

	_, result = getTypeAndResultFromTasks([]devops.CICDTask{successDomainTask})
	assert.Equal(t, `SUCCESS`, result)

	_, result = getTypeAndResultFromTasks([]devops.CICDTask{successDomainTask, successDomainTask})
	assert.Equal(t, `SUCCESS`, result)

	_, result = getTypeAndResultFromTasks([]devops.CICDTask{successDomainTask, successDomainTask, failureDomainTask})
	assert.Equal(t, `FAILURE`, result)

	_, result = getTypeAndResultFromTasks([]devops.CICDTask{successDomainTask, successDomainTask, abortDomainTask})
	assert.Equal(t, `ABORT`, result)

	_, result = getTypeAndResultFromTasks([]devops.CICDTask{failureDomainTask, failureDomainTask, abortDomainTask})
	assert.Equal(t, `ABORT`, result)
}
