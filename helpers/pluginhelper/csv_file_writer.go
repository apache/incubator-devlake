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
	"os"
)

// CsvFileWriter make writer for saving csv file easier, it write tuple to csv file
//
// Example CSV format (exported by dbeaver):
//
//   "id","name","json","created_at"
//   123,"foobar","{""url"": ""https://example.com""}","2022-05-05 09:56:43.438000000"
//
type CsvFileWriter struct {
	file   *os.File
	writer *csv.Writer
	fields []string
}

// NewCsvFileWriter create a `*CsvFileWriter` based on path to saving csv file
func NewCsvFileWriter(csvPath string, fields []string) *CsvFileWriter {
	// open csv file
	csvFile, err := os.Create(csvPath)
	if err != nil {
		panic(err)
	}
	csvWriter := csv.NewWriter(csvFile)
	// write field names
	err = csvWriter.Write(fields)
	if err != nil {
		panic(err)
	}
	csvWriter.Flush()
	if err != nil {
		panic(err)
	}
	return &CsvFileWriter{
		file:   csvFile,
		writer: csvWriter,
		fields: fields,
	}
}

// Close releases resource
func (ci *CsvFileWriter) Close() {
	ci.writer.Flush()
	err := ci.file.Close()
	if err != nil {
		panic(err)
	}
}

// Write the values into csv
func (ci *CsvFileWriter) Write(values []string) {
	err := ci.writer.Write(values)
	if err != nil {
		panic(err)
	}
}

// Flush the wrote data into file physically
func (ci *CsvFileWriter) Flush() {
	ci.writer.Flush()
}
