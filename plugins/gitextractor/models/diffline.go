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

import git "github.com/libgit2/git2go/v33"

type Difflines []git.DiffLine

// some essential functions for custom sort
func (difflines Difflines) Len() int {
	return len(difflines)
}

func (difflines Difflines) Less(i, j int) bool {
	return difflines[i].OldLineno > difflines[j].OldLineno
}

func (difflines Difflines) Swap(i, j int) {
	temp := difflines[i]
	difflines[i] = difflines[j]
	difflines[j] = temp
}
