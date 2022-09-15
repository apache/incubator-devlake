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
	ciDomainTask := devops.CICDTask{Type: `CI`}
	cdDomainTask := devops.CICDTask{Type: `CD`}
	cicdDomainTask := devops.CICDTask{Type: `CI/CD`}

	abortDomainTask := devops.CICDTask{Result: `ABORT`}
	failureDomainTask := devops.CICDTask{Result: `FAILURE`}
	successDomainTask := devops.CICDTask{Result: `SUCCESS`}

	hasCi, hasCd, _ := getTypeAndResultFromTasks([]devops.CICDTask{ciDomainTask})
	assert.True(t, hasCi)
	assert.False(t, hasCd)

	hasCi, hasCd, _ = getTypeAndResultFromTasks([]devops.CICDTask{cdDomainTask, cdDomainTask, cdDomainTask})
	assert.False(t, hasCi)
	assert.True(t, hasCd)

	hasCi, hasCd, _ = getTypeAndResultFromTasks([]devops.CICDTask{ciDomainTask, cicdDomainTask, cdDomainTask})
	assert.True(t, hasCi)
	assert.True(t, hasCd)

	hasCi, hasCd, _ = getTypeAndResultFromTasks([]devops.CICDTask{ciDomainTask, cdDomainTask, cdDomainTask})
	assert.True(t, hasCi)
	assert.True(t, hasCd)

	_, _, result := getTypeAndResultFromTasks([]devops.CICDTask{abortDomainTask, failureDomainTask, successDomainTask})
	assert.Equal(t, `ABORT`, result)

	_, _, result = getTypeAndResultFromTasks([]devops.CICDTask{failureDomainTask, successDomainTask})
	assert.Equal(t, `FAILURE`, result)

	_, _, result = getTypeAndResultFromTasks([]devops.CICDTask{successDomainTask})
	assert.Equal(t, `SUCCESS`, result)

	_, _, result = getTypeAndResultFromTasks([]devops.CICDTask{successDomainTask, successDomainTask})
	assert.Equal(t, `SUCCESS`, result)

	_, _, result = getTypeAndResultFromTasks([]devops.CICDTask{successDomainTask, successDomainTask, failureDomainTask})
	assert.Equal(t, `FAILURE`, result)

	_, _, result = getTypeAndResultFromTasks([]devops.CICDTask{successDomainTask, successDomainTask, abortDomainTask})
	assert.Equal(t, `ABORT`, result)

	_, _, result = getTypeAndResultFromTasks([]devops.CICDTask{failureDomainTask, failureDomainTask, abortDomainTask})
	assert.Equal(t, `ABORT`, result)
}
