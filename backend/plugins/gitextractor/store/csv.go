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

package store

import (
	"encoding/csv"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"os"
	"path/filepath"
	"reflect"
)

type csvWriter struct {
	f *os.File
	w *csv.Writer
}

func newCsvWriter(path string, v interface{}) (*csvWriter, errors.Error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, errors.Convert(err)
	}
	// declare UTF-8 encoding
	_, err = f.WriteString("\xEF\xBB\xBF")
	if err != nil {
		return nil, errors.Convert(err)
	}
	w := csv.NewWriter(f)
	value := reflect.Indirect(reflect.ValueOf(v))
	var header []string
	for i := 0; i < value.NumField(); i++ {
		if value.Type().Field(i).Anonymous {
			continue
		}
		header = append(header, value.Type().Field(i).Name)
	}
	err = w.Write(header)
	if err != nil {
		return nil, errors.Convert(err)
	}
	return &csvWriter{f: f, w: w}, nil
}

func (w *csvWriter) Write(item interface{}) errors.Error {
	v := reflect.Indirect(reflect.ValueOf(item))
	n := v.NumField()
	record := make([]string, 0, n)
	for i := 0; i < n; i++ {
		if v.Type().Field(i).Anonymous {
			continue
		}
		record = append(record, fmt.Sprint(v.Field(i).Interface()))
	}
	return errors.Convert(w.w.Write(record))
}

func (w *csvWriter) Close() errors.Error {
	w.w.Flush()
	return errors.Convert(w.f.Close())
}

type CsvStore struct {
	dir                       string
	repoCommitWriter          *csvWriter
	commitWriter              *csvWriter
	refWriter                 *csvWriter
	commitFileWriter          *csvWriter
	commitParentWriter        *csvWriter
	commitFileComponentWriter *csvWriter
	commitLineChangeWriter    *csvWriter
	snapshotWriter            *csvWriter
}

func NewCsvStore(dir string) (*CsvStore, errors.Error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return nil, errors.Convert(err)
		}
	}
	var err error
	s := &CsvStore{dir: dir}
	s.repoCommitWriter, err = newCsvWriter(filepath.Join(dir, "repo_commits.csv"), code.RepoCommit{})
	if err != nil {
		return nil, errors.Convert(err)
	}
	s.commitWriter, err = newCsvWriter(filepath.Join(dir, "commits.csv"), code.Commit{})
	if err != nil {
		return nil, errors.Convert(err)
	}
	s.refWriter, err = newCsvWriter(filepath.Join(dir, "refs.csv"), code.Ref{})
	if err != nil {
		return nil, errors.Convert(err)
	}
	s.commitFileWriter, err = newCsvWriter(filepath.Join(dir, "commit_files.csv"), code.CommitFile{})
	if err != nil {
		return nil, errors.Convert(err)
	}
	s.commitParentWriter, err = newCsvWriter(filepath.Join(dir, "commit_parents.csv"), code.CommitParent{})
	if err != nil {
		return nil, errors.Convert(err)
	}
	s.commitFileComponentWriter, err = newCsvWriter(filepath.Join(dir, "commit_file_components.csv"), code.CommitFileComponent{})
	if err != nil {
		return nil, errors.Convert(err)
	}
	s.commitLineChangeWriter, err = newCsvWriter(filepath.Join(dir, "commit_line_changes.csv"), code.CommitLineChange{})
	if err != nil {
		return nil, errors.Convert(err)
	}
	s.snapshotWriter, err = newCsvWriter(filepath.Join(dir, "repo_snapshot.csv"), code.RepoSnapshot{})
	if err != nil {
		return nil, errors.Convert(err)
	}
	return s, nil
}

func (c *CsvStore) RepoCommits(repoCommit *code.RepoCommit) errors.Error {
	return c.repoCommitWriter.Write(repoCommit)
}

func (c *CsvStore) Commits(commit *code.Commit) errors.Error {
	return c.commitWriter.Write(commit)
}

func (c *CsvStore) Refs(ref *code.Ref) errors.Error {
	return c.refWriter.Write(ref)
}

func (c *CsvStore) CommitFiles(file *code.CommitFile) errors.Error {
	return c.commitFileWriter.Write(file)
}

func (c *CsvStore) CommitFileComponents(commitFileComponent *code.CommitFileComponent) errors.Error {
	return c.commitFileComponentWriter.Write(commitFileComponent)
}

func (c *CsvStore) CommitLineChange(commitLineChange *code.CommitLineChange) errors.Error {
	return c.commitLineChangeWriter.Write(commitLineChange)
}

func (c *CsvStore) RepoSnapshot(ss *code.RepoSnapshot) errors.Error {
	return c.snapshotWriter.Write(ss)
}

func (c *CsvStore) CommitParents(pp []*code.CommitParent) errors.Error {
	var err error
	for _, p := range pp {
		err = c.commitParentWriter.Write(p)
		if err != nil {
			return errors.Convert(err)
		}
	}
	return nil
}

func (c *CsvStore) Close() errors.Error {
	if c.repoCommitWriter != nil {
		c.repoCommitWriter.Close()
	}
	if c.commitWriter != nil {
		c.commitWriter.Close()
	}
	if c.refWriter != nil {
		c.refWriter.Close()
	}
	if c.commitFileWriter != nil {
		c.commitFileWriter.Close()
	}
	if c.commitParentWriter != nil {
		c.commitParentWriter.Close()
	}
	if c.snapshotWriter != nil {
		c.snapshotWriter.Close()
	}
	if c.commitFileComponentWriter != nil {
		c.commitFileComponentWriter.Close()
	}
	if c.commitLineChangeWriter != nil {
		c.commitLineChangeWriter.Close()
	}
	return nil
}
