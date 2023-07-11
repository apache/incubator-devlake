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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type iteratorConcator struct {
	index     int
	iterators []api.Iterator
}

func newIteratorConcator(iterators ...api.Iterator) *iteratorConcator {
	return &iteratorConcator{iterators: iterators}
}

func (w *iteratorConcator) HasNext() bool {
	for w.index < len(w.iterators) {
		if w.iterators[w.index].HasNext() {
			return true
		}
		w.index++
	}
	return false
}

func (w *iteratorConcator) Fetch() (interface{}, errors.Error) {
	if w.index >= len(w.iterators) {
		return nil, errors.Default.New("index out of range")
	}
	return w.iterators[w.index].Fetch()
}

func (w *iteratorConcator) Close() errors.Error {
	for _, iterator := range w.iterators {
		iterator.Close()
	}
	return nil
}

type iteratorWrapper struct {
	original    api.Iterator
	wrapperFunc func(interface{}) interface{}
}

func newIteratorWrapper(original api.Iterator, wrapperFunc func(interface{}) interface{}) *iteratorWrapper {
	return &iteratorWrapper{original: original, wrapperFunc: wrapperFunc}
}

func (w *iteratorWrapper) HasNext() bool {
	return w.original.HasNext()
}

func (w *iteratorWrapper) Fetch() (interface{}, errors.Error) {
	data, err := w.original.Fetch()
	if err != nil {
		return nil, err
	}
	return w.wrapperFunc(data), nil
}

func (w *iteratorWrapper) Close() errors.Error {
	return w.original.Close()
}

type iteratorFromSlice struct {
	index int
	data  []interface{}
}

func newIteratorFromSlice(data []interface{}) *iteratorFromSlice {
	return &iteratorFromSlice{data: data}
}

func (i *iteratorFromSlice) HasNext() bool {
	return i.index < len(i.data)
}

func (i *iteratorFromSlice) Fetch() (interface{}, errors.Error) {
	data := i.data[i.index]
	i.index++
	return data, nil
}

func (i *iteratorFromSlice) Close() errors.Error {
	i.index = len(i.data) - 1
	return nil
}
