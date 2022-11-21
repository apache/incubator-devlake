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
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

// ErrResponseNotExist represents the requested uri has never been cached
var ErrResponseNotExist = errors.New("resopnse not exist")

// ResponseCache stores and provides the status/headers/body of a response
type ResponseCache interface {
	Get(key string) (int, http.Header, []byte, error)
	Set(key string, status int, headers http.Header, body []byte) error
}

// ResponseDiskCache caches response to specified folder
type ResponseDiskCache struct {
	folder string
}

// Get status, headers, body for a specific key
func (d *ResponseDiskCache) Get(key string) (int, http.Header, []byte, error) {
	// make sure file exists and it is in fact a file
	fp := path.Join(d.folder, key)
	if fs, err := os.Stat(fp); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil, nil, ErrResponseNotExist
		}
		panic(err)
	} else if fs.IsDir() {
		panic(fmt.Errorf("%v is a folder", fp))
	}

	// read file line by line to get status,header and body

	file, err := os.Open(fp)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	readLine := func() string {
		buf, err := reader.ReadBytes('\r')
		if err != nil {
			panic(err)
		}
		nl, err := reader.ReadByte()
		if err != nil {
			panic(err)
		}
		if nl != '\n' {
			panic(fmt.Errorf("expecting \\r followed by \\n"))
		}
		return string(buf[:len(buf)-1])
	}
	status, err := strconv.Atoi(readLine())
	if err != nil {
		panic(err)
	}
	headers := make(http.Header)
	for {
		header := readLine()
		if header == "" {
			break
		}
		idx := strings.Index(header, ": ")
		if idx < 1 {
			panic(fmt.Sprintf("unexpected colon position for header %v", header))
		}
		name := header[:idx]
		value := header[idx+2:]
		headers.Add(name, value)
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return status, headers, body[:len(body)-2], nil
}

// Set a response cache by key
func (d *ResponseDiskCache) Set(key string, status int, headers http.Header, body []byte) error {
	fp := path.Join(d.folder, key)
	file, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writeLine := func(line []byte) {
		buf := append(line, '\r', '\n')
		n, err := file.Write(buf)
		if err != nil {
			panic(err)
		}
		if n != len(buf) {
			panic(fmt.Errorf("write cache failed, expected %v bytes written, got %v", len(buf), n))
		}
	}

	writeLine([]byte(fmt.Sprintf("%v", status)))
	for name, values := range headers {
		for _, value := range values {
			writeLine([]byte(fmt.Sprintf("%v: %v", name, value)))
		}
	}
	writeLine(nil)

	writeLine(body)
	return nil
}

// NewDiskCache creates a new disk cache for response storage
func NewDiskCache(folder string) ResponseCache {
	if fs, err := os.Stat(folder); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.Mkdir(folder, 0755)
			if err != nil {
				panic(fmt.Sprintf("failed to create cache folder %v", folder))
			}
		}
	} else if !fs.IsDir() {
		panic(fmt.Sprintf("%v is not a folder", folder))
	}
	return &ResponseDiskCache{folder: folder}
}
