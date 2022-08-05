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
)

type pathStats struct {
	info     os.FileInfo
	path     string
	filename string
	isNested bool
}

// CreateArchive creates a tar archive and writes the files/directories associated with the `sourcePaths` to it.
// source adapted from https://golangdocs.com/tar-gzip-in-golang.
func CreateArchive(archivePath string, sourcePaths ...string) error {
	tarfile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("error creating archive %s: %v", archivePath, err)
	}
	defer tarfile.Close()
	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()
	for _, sourcePath := range sourcePaths {
		var source *pathStats
		source, err = getPathStats(sourcePath)
		if err != nil {
			return err
		}
		err = filepath.Walk(sourcePath, func(subPath string, subPathInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(subPathInfo, subPathInfo.Name())
			if err != nil {
				return err
			}
			header.Name = subPath

			// should skip if the src is a nested file, and this is a file that does not exactly match that
			shouldSkip := !subPathInfo.IsDir() && source.isNested &&
				source.filename != subPathInfo.Name() &&
				filepath.Dir(source.filename) == subPathInfo.Name()
			if shouldSkip {
				return nil
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
