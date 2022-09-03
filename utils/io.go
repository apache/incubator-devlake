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
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/viant/afs"
	"os"
	"path/filepath"
	"strings"
)

// fs abstract filesystem interface singleton instance
var fs = afs.New()

// CreateZipArchive creates a zip archive and writes the files/directories associated with the `sourcePaths` to it.
// If a sourcePath directory ends with /*, then its contents are copied over, but not the directory itself
func CreateZipArchive(archivePath string, sourcePaths ...string) errors.Error {
	return createArchive("zip", archivePath, sourcePaths...)
}

// CreateGZipArchive creates a tar archive, compresses it with gzip and writes the files/directories associated with the `sourcePaths` to it.
// If a sourcePath directory ends with /*, then its contents are copied over, but not the directory itself
func CreateGZipArchive(archivePath string, sourcePaths ...string) errors.Error {
	err := createArchive("tar", archivePath, sourcePaths...)
	if err != nil {
		return err
	}
	// now gzip it
	err = toGzip(archivePath)
	if err != nil {
		return errors.Default.Wrap(err, "error compressing archive to gzip")
	}
	return nil
}

func createArchive(archiveType string, archivePath string, sourcePaths ...string) errors.Error {
	for _, sourcePath := range sourcePaths {
		relativeCopy := false
		if strings.HasSuffix(sourcePath, "/*") {
			sourcePath = sourcePath[0 : len(sourcePath)-2]
			relativeCopy = true
		}
		srcPathAbs, err := filepath.Abs(sourcePath)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("error getting absolute path of %s", sourcePath))
		}
		archivePathAbs, err := filepath.Abs(archivePath)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("error getting absolute path of %s", archivePath))
		}
		srcInfo, err := os.Stat(srcPathAbs)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("error getting stats of path %s", srcPathAbs))
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
			return errors.Default.Wrap(err, "error trying to copy data to archive")
		}
	}
	return nil
}

func copyContentsToArchive(archiveType string, absSourcePath string, absArchivePath string) errors.Error {
	var files []os.DirEntry
	files, err := os.ReadDir(absSourcePath)
	if err != nil {
		return errors.Convert(err)
	}
	for _, desPath := range files {
		archiveDest := desPath.Name()
		src := fmt.Sprintf("%s/%s", absSourcePath, archiveDest)
		err = copyToArchive(archiveType, src, absArchivePath, archiveDest)
		if err != nil {
			return errors.Convert(err)
		}
	}
	return nil
}

func copyToArchive(archiveType string, absSourcePath string, absArchivePath string, archiveDest string) errors.Error {
	src := fmt.Sprintf("file://%s", absSourcePath)
	dst := fmt.Sprintf("file:%s/%s:///%s", absArchivePath, archiveType, archiveDest)
	return errors.Convert(fs.Copy(context.Background(), src, dst))
}

func toGzip(archivePath string) errors.Error {
	info, _ := os.Stat(archivePath)
	b, err := os.ReadFile(archivePath)
	if err != nil {
		return errors.Convert(err)
	}
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err = w.Write(b)
	_ = w.Close()
	if err != nil {
		return errors.Convert(err)
	}
	return errors.Convert(os.WriteFile(archivePath, buf.Bytes(), info.Mode().Perm()))
}
