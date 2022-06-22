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
	"strings"

	"github.com/apache/incubator-devlake/plugins/core"
)

type DomainIdGenerator struct {
	prefix  string
	pkNames []string
	pkTypes []reflect.Type
}

type WildCard string

const WILDCARD WildCard = "%"

var wildcardType = reflect.TypeOf(WILDCARD)

func walkFields(t reflect.Type, pkNames *[]string, pkTypes *[]reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			walkFields(field.Type, pkNames, pkTypes)
		} else {
			gormTag := field.Tag.Get("gorm")

			// TODO: regex?
			if gormTag != "" && strings.Contains(strings.ToLower(gormTag), "primarykey") {
				*pkNames = append(*pkNames, field.Name)
				*pkTypes = append(*pkTypes, field.Type)
			}
		}
	}
}

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

	// find out all primkary keys and their types
	pkNames := make([]string, 0, 1)
	pkTypes := make([]reflect.Type, 0, 1)

	walkFields(t, &pkNames, &pkTypes)

	if len(pkNames) == 0 {
		panic(fmt.Errorf("no primary key found for %s:%s", pluginName, structName))
	}

	return &DomainIdGenerator{
		prefix:  fmt.Sprintf("%s:%s", pluginName, structName),
		pkNames: pkNames,
		pkTypes: pkTypes,
	}
}

func (g *DomainIdGenerator) Generate(pkValues ...interface{}) string {
	id := g.prefix
	isMatch := false

	for _, pkValue := range pkValues {
		// append pk
		id += ":" + fmt.Sprintf("%v", pkValue)
		// type checking
		pkValueType := reflect.TypeOf(pkValue)
		if pkValueType == wildcardType {
			break
		}

		for _, pkType := range g.pkTypes {
			if pkValueType == pkType {
				isMatch = true
			}
		}

		if !isMatch {
			panic(fmt.Errorf("primary key type [%s] does not match: %s",
				pkValueType,
				g.pkTypes,
			))
		}

	}
	return id
}
