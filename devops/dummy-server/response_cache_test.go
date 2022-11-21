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

package main

import (
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiskResponseCache(t *testing.T) {
	cacheDir := path.Join(os.TempDir(), "cacheDir")
	assert.Nil(t, os.RemoveAll(cacheDir))
	diskCache := NewDiskCache(cacheDir)
	assert.Nil(t, diskCache.Set("hello", 200, http.Header{"foo": []string{"bar"}}, []byte("world")))
	status, headers, body, err := diskCache.Get("hello")
	assert.Nil(t, err)
	assert.Equal(t, 200, status)
	assert.Equal(t, http.Header{"Foo": []string{"bar"}}, headers)
	assert.Equal(t, "world", string(body))
}
