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

package errors

import "sync"

// Wraps the native sync map in generic form. Consider moving this to another package for broader use later
type syncMap[K any, V any] struct {
	m *sync.Map
}

func newSyncMap[K any, V any]() *syncMap[K, V] {
	return &syncMap[K, V]{
		m: new(sync.Map),
	}
}

func (sm *syncMap[K, V]) Store(key K, val V) {
	sm.m.Store(key, val)
}

func (sm *syncMap[K, V]) Load(key K) (V, bool) {
	var v V
	val, ok := sm.m.Load(key)
	if ok {
		v = val.(V)
	}
	return v, ok
}
