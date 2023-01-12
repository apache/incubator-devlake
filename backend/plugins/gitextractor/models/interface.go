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

package models

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
)

type Store interface {
	RepoCommits(repoCommit *code.RepoCommit) errors.Error
	Commits(commit *code.Commit) errors.Error
	Refs(ref *code.Ref) errors.Error
	CommitFiles(file *code.CommitFile) errors.Error
	CommitParents(pp []*code.CommitParent) errors.Error
	CommitFileComponents(commitFileComponent *code.CommitFileComponent) errors.Error
	CommitLineChange(commitLineChange *code.CommitLineChange) errors.Error
	RepoSnapshot(snapshot *code.RepoSnapshot) errors.Error
	Close() errors.Error
}
