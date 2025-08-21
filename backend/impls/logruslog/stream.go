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

package logruslog

import (
	"github.com/apache/incubator-devlake/core/errors"
	"io"
	"os"
	"path/filepath"
)

func GetFileStream(path string) (io.Writer, errors.Error) {
	if path == "" {
		return os.Stdout, nil
	}
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return nil, errors.Convert(err)
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, errors.Convert(err)
	}
	return io.MultiWriter(os.Stdout, file), nil
}
