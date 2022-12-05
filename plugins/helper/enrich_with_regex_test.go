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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegexEnricher_GetEnrichResult(t *testing.T) {
	re := NewRegexEnricher()
	pattern1 := `deploy.*`
	pattern2 := "production"
	err := re.AddRegexp(pattern1, pattern2)
	assert.Nil(t, err)
	res1 := re.GetEnrichResult(pattern1, `deployToWin`, `deployment`)
	assert.Equal(t, "deployment", res1)
	res2 := re.GetEnrichResult(pattern1, `deplo1`, `deployment`)
	assert.Equal(t, "", res2)

	res3 := re.GetEnrichResult(pattern2, `production`, `product`)
	assert.Equal(t, "product", res3)
	res4 := re.GetEnrichResult(pattern2, `producti1n`, `product`)
	assert.Equal(t, "", res4)
}
