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

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ PluginMeta = (*Foo)(nil)
var _ PluginMeta = (*Bar)(nil)

type Foo string

func (f *Foo) Description() string {
	return "foo"
}

func (f *Foo) RootPkgPath() string {
	return "path/to/foo"
}

type Bar string

func (b *Bar) Description() string {
	return "foo"
}

func (b *Bar) RootPkgPath() string {
	return "path/to/bar"
}

func TestHub(t *testing.T) {
	var foo Foo
	assert.Nil(t, RegisterPlugin("foo", &foo))
	var bar Bar
	assert.Nil(t, RegisterPlugin("bar", &bar))

	f, _ := GetPlugin("foo")
	assert.Equal(t, &foo, f)

	fn, _ := FindPluginNameBySubPkgPath("path/to/foo/models")
	assert.Equal(t, fn, "foo")

	b, _ := GetPlugin("bar")
	assert.Equal(t, &bar, b)

	bn, _ := FindPluginNameBySubPkgPath("path/to/bar/models")
	assert.Equal(t, bn, "bar")
}
