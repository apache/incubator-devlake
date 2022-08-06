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

package utils

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
import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type pathStats struct {
	info     os.FileInfo
	path     string
	filename string
	isNested bool
}

// CreateArchive creates a tar archive and writes the files/directories associated with the 'sourcePaths' to it.
//
// 'archivePath': the path on the filesystem where the archive file will be located.
//
// 'relativeCopy': if true, it'll only copy the contents inside a sourcePath (if sourcePath is a file, then just the file), otherwise the sourcePath is recursively copied.
//
// 'sourcePaths': each sourcePath may be a file or a directory. A file nested in a directory (i.e. a/b/x.log) will be copied along with its parent folders unless relativeCopy is true.
func CreateArchive(archivePath string, relativeCopy bool, sourcePaths ...string) error {
	err := os.MkdirAll(filepath.Dir(archivePath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating path for %s: %v", archivePath, err)
	}
	tarfile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("error creating archive %s: %v", archivePath, err)
	}
	defer tarfile.Close()
	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()
	for _, sourcePath := range sourcePaths {
		var sourceStats *pathStats
		sourceStats, err = getPathStats(sourcePath)
		if err != nil {
			return err
		}
		err = filepath.Walk(sourcePath, func(subPath string, subPathInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// if we're copying the contents of a directory 'relatively' then skip until we reach that directory
			if relativeCopy && !sourceStats.isNested && isParentTo(subPath, sourcePath) {
				return nil
			}
			// should skip if the src is a nested file, and this is a file that does not exactly match it
			if !subPathInfo.IsDir() && sourceStats.isNested && sourceStats.filename != subPathInfo.Name() &&
				filepath.Dir(sourceStats.filename) == subPathInfo.Name() {
				return nil
			}
			header, err := tar.FileInfoHeader(subPathInfo, subPathInfo.Name())
			if err != nil {
				return err
			}
			if relativeCopy {
				// adjust the header to the correct relative path
				header.Name = removeParent(sourcePath, subPath)
			} else {
				header.Name = subPath
			}
			if err = tarball.WriteHeader(header); err != nil {
				return err
			}
			if subPathInfo.IsDir() {
				return nil
			}
			file, err := os.Open(subPath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
		if err != nil {
			return fmt.Errorf("error copying %s into archive %s: %v", sourcePath, archivePath, err)
		}
	}
	return nil
}

func getPathStats(sourcePath string) (*pathStats, error) {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("error getting metadata about %s", sourcePath)
	}
	filename := filepath.Base(sourcePath)
	if !sourceInfo.IsDir() && filepath.Dir(sourcePath) != "." {
		sourcePath = filepath.Dir(sourcePath)
		sourceInfo, _ = os.Stat(sourcePath)
		return &pathStats{
			info:     sourceInfo,
			path:     sourcePath,
			filename: filename,
			isNested: true,
		}, nil
	}
	return &pathStats{
		info:     sourceInfo,
		path:     sourcePath,
		filename: filename,
		isNested: false,
	}, nil
}

// example: (a/b, a/b/c) -> true
func isParentTo(parent string, path string) bool {
	p := path
	for p != "." {
		if p == parent {
			return true
		}
		p = filepath.Dir(p)
	}
	return false
}

// example (a/b, a/b/c) -> c
func removeParent(parent string, child string) string {
	s := strings.TrimPrefix(strings.TrimPrefix(child, parent), "/")
	if s == "" {
		return filepath.Base(child)
	}
	return s
}
