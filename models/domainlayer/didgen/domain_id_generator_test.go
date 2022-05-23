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

package didgen

import (
	"context"
	"testing"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type FooPlugin string

func (f *FooPlugin) Description() string {
	return "foo"
}

func (f *FooPlugin) Init() {
}

func (f *FooPlugin) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	return nil
}

func (f *FooPlugin) RootPkgPath() string {
	return "github.com/apache/incubator-devlake"
}

func (f *FooPlugin) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return make(map[string]map[string]core.ApiResourceHandler)
}

type FooModel struct {
	gorm.Model
}

func TestOriginKeyGenerator(t *testing.T) {
	var foo FooPlugin
	assert.Nil(t, core.RegisterPlugin("fooplugin", &foo))

	g := NewDomainIdGenerator(&FooModel{})
	assert.Equal(t, g.prefix, "fooplugin:FooModel")

	originKey := g.Generate(uint(2))
	assert.Equal(t, "fooplugin:FooModel:2", originKey)

	assert.Panics(t, func() {
		NewDomainIdGenerator(&foo)
	})

	assert.Panics(t, func() {
		g.Generate("asdf")
	})
}
