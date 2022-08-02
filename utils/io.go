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
	"context"
	"fmt"
	"github.com/viant/afs"
	"path/filepath"
)

// fs abstract filesystem interface singleton instance
var fs = afs.New()

// CreateArchive creates a tar archive and writes the files/directories associated with the `targetPaths` to it.
func CreateArchive(archivePath string, targetPaths ...string) error {
	var err error
	for _, targetPath := range targetPaths {
		targetPath, err = filepath.Abs(targetPath)
		if err != nil {
			return err
		}
		archivePath, err = filepath.Abs(archivePath)
		if err != nil {
			return err
		}
		err = fs.Copy(context.Background(), fmt.Sprintf("file://%s", targetPath), fmt.Sprintf("file:%s/tar:///", archivePath))
		if err != nil {
			return err
		}
	}
	return nil
}
