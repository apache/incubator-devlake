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
	"strings"
)

// fs abstract filesystem interface singleton instance
var fs = afs.New()

// CreateZipArchive creates a zip archive and writes the files/directories associated with the `targetPaths` to it.
// If a sourcePath directory ends with /*, then its contents are copied over, but not the directory itself
func CreateZipArchive(archivePath string, sourcePaths ...string) error {
	return createArchive("zip", archivePath, sourcePaths...)
}

func createArchive(archiveType string, archivePath string, sourcePaths ...string) error {
	for _, sourcePath := range sourcePaths {
		relativeCopy := false
		if strings.HasSuffix(sourcePath, "/*") {
			sourcePath = sourcePath[0 : len(sourcePath)-2]
			relativeCopy = true
		}
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
		if relativeCopy && srcInfo.IsDir() {
			err = copyContentsToArchive(archiveType, srcPathAbs, archivePathAbs)
		} else {
			if relativeCopy {
				sourcePath = srcInfo.Name()
			}
			// directly copies over the src path as the dest path in the archive (can be improved later if needed to guard against abs paths)
			err = copyToArchive(archiveType, srcPathAbs, archivePathAbs, sourcePath)
		}
		if err != nil {
			return fmt.Errorf("error trying to copy data to archive: %v", err)
		}
	}
	return nil
}

func copyContentsToArchive(archiveType string, absSourcePath string, absArchivePath string) error {
	var files []os.DirEntry
	files, err := os.ReadDir(absSourcePath)
	if err != nil {
		return err
	}
	for _, desPath := range files {
		archiveDest := desPath.Name()
		src := fmt.Sprintf("%s/%s", absSourcePath, archiveDest)
		err = copyToArchive(archiveType, src, absArchivePath, archiveDest)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyToArchive(archiveType string, absSourcePath string, absArchivePath string, archiveDest string) error {
	src := fmt.Sprintf("file://%s", absSourcePath)
	dst := fmt.Sprintf("file:%s/%s:///%s", absArchivePath, archiveType, archiveDest)
	return fs.Copy(context.Background(), src, dst)
}
