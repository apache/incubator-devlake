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
	"fmt"
	"reflect"

	"github.com/apache/incubator-devlake/impl/dalgorm"
	"github.com/apache/incubator-devlake/plugins/core"
)

type DomainIdGenerator struct {
	prefix string
	pk     []reflect.StructField
}

type WildCard string

const WILDCARD WildCard = "%"

var wildcardType = reflect.TypeOf(WILDCARD)

func NewDomainIdGenerator(entityPtr interface{}) *DomainIdGenerator {
	v := reflect.ValueOf(entityPtr)
	if v.Kind() != reflect.Ptr {
		panic("entityPtr is not a pointer")
	}
	t := reflect.Indirect(v).Type()

	// find out which plugin holds the entity
	pluginName, err := core.FindPluginNameBySubPkgPath(t.PkgPath())
	if err != nil {
		panic(err)
	}
	// find out entity type name
	structName := t.Name()

	dal := &dalgorm.Dalgorm{}
	pk := dal.GetPrimarykeyFields(t)

	if len(pk) == 0 {
		panic(fmt.Errorf("no primary key found for %s:%s", pluginName, structName))
	}

	return &DomainIdGenerator{
		prefix: fmt.Sprintf("%s:%s", pluginName, structName),
		pk:     pk,
	}
}

func (g *DomainIdGenerator) Generate(pkValues ...interface{}) string {
	id := g.prefix
	for i, pkValue := range pkValues {
		// append pk
		id += ":" + fmt.Sprintf("%v", pkValue)
		// type checking
		pkValueType := reflect.TypeOf(pkValue)
		if pkValueType == wildcardType {
			break
		} else if pkValueType != g.pk[i].Type {
			panic(fmt.Errorf("primary key type does not match: %s is %s type, and it should be %s type",
				g.pk[i].Name,
				pkValueType.Name(),
				g.pk[i].Type.Name(),
			))
		}
	}
	return id
}
