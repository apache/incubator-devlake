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

import (
	"context"
	"fmt"
	"github.com/viant/afs"
	"os"
	"path/filepath"
)

// fs abstract filesystem interface singleton instance
var fs = afs.New()

// CreateArchive creates a tar archive and writes the files/directories associated with the `targetPaths` to it.
// 'relativeCopy' = true will copy the contents inside each sourcePath (directory) over. If the sourcePath is a file, it is directly copied over.
func CreateArchive(archivePath string, relativeCopy bool, sourcePaths ...string) error {
	for _, sourcePath := range sourcePaths {
		srcPathAbs, err := filepath.Abs(sourcePath)
		if err != nil {
			return fmt.Errorf("error getting absolute path of %s: %v", sourcePaths, err)
		}
		archivePathAbs, err := filepath.Abs(archivePath)
		if err != nil {
			return fmt.Errorf("error getting absolute path of %s: %v", archivePath, err)
		}
		srcInfo, err := os.Stat(srcPathAbs)
		if err != nil {
			return fmt.Errorf("error getting stats of path %s: %v", srcPathAbs, err)
		}
		if relativeCopy {
			if srcInfo.IsDir() {
				err = copyContentsToArchive(archivePathAbs, archivePathAbs)
			} else {
				err = copyToArchive(srcPathAbs, archivePathAbs, srcInfo.Name())
			}
		} else {
			err = copyToArchive(srcPathAbs, archivePathAbs, sourcePath)
		}
		if err != nil {
			return fmt.Errorf("error trying to copy data to archive: %v", err)
		}
	}
	return nil
}

func copyContentsToArchive(absSourcePath string, absArchivePath string) error {
	var files []os.DirEntry
	files, err := os.ReadDir(absSourcePath)
	if err != nil {
		return err
	}
	for _, file := range files {
		err = fs.Copy(context.Background(), fmt.Sprintf("file://%s/%s", absSourcePath, file.Name()), fmt.Sprintf("file:%s/tar:///%s", absArchivePath, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func copyToArchive(absSourcePath string, absArchivePath string, filename string) error {
	return fs.Copy(context.Background(), fmt.Sprintf("file://%s", absSourcePath), fmt.Sprintf("file:%s/tar:///%s", absArchivePath, filename))
}
