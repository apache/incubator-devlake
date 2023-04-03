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
	"github.com/apache/incubator-devlake/core/errors"
	"os"
	"path/filepath"
)

var (
	ProjectRoot = ""
	Shell       = ""
)

func init() {
	Shell = "/bin/sh"
	var err errors.Error
	ProjectRoot, err = NormalizeBaseDirectory()
	if err != nil {
		panic(err.Error())
	}
}

func NormalizeBaseDirectory() (string, errors.Error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", errors.Convert(err)
	}
	for {
		dir := filepath.Base(pwd)
		if dir == "" {
			return "", errors.Default.New("base repo directory not found")
		}
		if dir == "backend" {
			break
		}
		pwd = filepath.Dir(pwd)
	}
	err = os.Chdir(pwd)
	return pwd, errors.Convert(err)
}
