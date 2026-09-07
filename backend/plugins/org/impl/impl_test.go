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

package impl

import (
	"testing"

	"github.com/apache/devlake/core/plugin"
	"github.com/stretchr/testify/assert"
)

// Route registration may evaluate ApiResources() before InitPlugins() has run
// (see https://github.com/apache/devlake/issues/9021). Handlers obtained from
// an uninitialized plugin must fail gracefully instead of panicking with a
// nil pointer dereference, and keep working once Init() runs later.
func TestApiResourcesBeforeInitFailsGracefully(t *testing.T) {
	p := &Org{}
	for resource, methods := range p.ApiResources() {
		for method, handler := range methods {
			assert.NotPanics(t, func() {
				out, err := handler(&plugin.ApiResourceInput{})
				assert.Nilf(t, out, "%s %s should return no output before Init", method, resource)
				assert.NotNilf(t, err, "%s %s should return an error before Init", method, resource)
			}, "%s %s must not panic before Init", method, resource)
		}
	}
}
