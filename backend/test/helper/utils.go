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
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apache/incubator-devlake/core/plugin"
)

// IsWSL FIXME
func IsWSL() bool {
	lines, err := readFile("/proc/version")
	if err != nil {
		return false
	}
	for _, line := range lines {
		l := strings.ToLower(line)
		if strings.Contains(l, "microsoft") {
			return true
		}
	}
	return false
}

// GetSubtaskNames FIXME
func GetSubtaskNames(metas ...plugin.SubTaskMeta) []string {
	var names []string
	for _, m := range metas {
		names = append(names, m.Name)
	}
	return names
}

// AddToPath FIXME
func AddToPath(newPaths ...string) {
	path := os.ExpandEnv("$PATH")
	for _, newPath := range newPaths {
		newPath, _ = filepath.Abs(newPath)
		path = fmt.Sprintf("%s:%s", newPath, path)
	}
	_ = os.Setenv("PATH", path)
}

func Val[T any](t T) *T {
	return &t
}

func Cast[T any](m any) T {
	j := ToJson(m)
	t := new(T)
	err := json.Unmarshal(j, t)
	if err != nil {
		panic(err)
	}
	return *t
}

func readFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
