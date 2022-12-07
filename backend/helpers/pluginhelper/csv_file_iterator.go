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

package pluginhelper

import (
	"encoding/csv"
	"io"
	"os"
)

// CsvFileIterator make iterating rows from csv file easier, it reads tuple from csv file and turn it into
// a `map[string]interface{}` for you.
//
// Example CSV format (exported by dbeaver):
//
//	"id","name","json","created_at"
//	123,"foobar","{""url"": ""https://example.com""}","2022-05-05 09:56:43.438000000"
type CsvFileIterator struct {
	file   *os.File
	reader *csv.Reader
	fields []string
	row    map[string]interface{}
}

// NewCsvFileIterator create a `*CsvFileIterator` based on path to csv file
func NewCsvFileIterator(csvPath string) *CsvFileIterator {
	// open csv file
	csvFile, err := os.Open(csvPath)
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(csvFile)
	// load field names
	fields, err := csvReader.Read()
	if err != nil {
		panic(err)
	}
	return &CsvFileIterator{
		file:   csvFile,
		reader: csvReader,
		fields: fields,
	}
}

// Close releases resource
func (ci *CsvFileIterator) Close() {
	err := ci.file.Close()
	if err != nil {
		panic(err)
	}
}

// HasNext returns a boolean to indicate whether there was any row to be `Fetch`
func (ci *CsvFileIterator) HasNext() bool {
	row, err := ci.reader.Read()
	if err == io.EOF {
		ci.row = nil
		return false
	}
	if err != nil {
		ci.row = nil
		panic(err)
	}
	// convert row tuple to map type, so gorm can insert data with it
	ci.row = make(map[string]interface{})
	for index, field := range ci.fields {
		ci.row[field] = row[index]
	}
	return true
}

// Fetch returns current row
func (ci *CsvFileIterator) Fetch() map[string]interface{} {
	return ci.row
}
