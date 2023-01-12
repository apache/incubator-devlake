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
	"container/list"
)

type FileBlame struct {
	Idx   int
	It    *list.Element
	Lines *list.List
}

// Walk to a specific index
func (fb *FileBlame) Walk(num int) {
	for fb.Idx < num && fb.It != fb.Lines.Back() {
		fb.Idx++
		fb.It = fb.It.Next()
	}
	for fb.Idx > num && fb.It != fb.Lines.Front() {
		fb.Idx--
		fb.It = fb.It.Prev()
	}
}

// Find an element with specific line number
func (fb *FileBlame) Find(num int) *list.Element {
	fb.Walk(num)
	if fb.Idx == num && fb.It != nil {
		return fb.It
	}
	return nil
}

// AddLine Add a line at a specific line num
func (fb *FileBlame) AddLine(num int, commit string) {
	fb.Walk(num)
	flag := false
	for fb.It == fb.Lines.Back() && fb.Idx < num {
		flag = true
		fb.It = fb.Lines.PushBack(nil)
		fb.Idx++

	}
	if fb.It == nil {
		fb.It = fb.Lines.PushBack(commit)
	} else if flag {
		fb.It.Value = commit
	} else {
		fb.It = fb.Lines.InsertBefore(commit, fb.It)
	}
}

// RemoveLine remove a line at num
func (fb *FileBlame) RemoveLine(num int) {
	fb.Walk(num)
	a := fb.It
	if fb.Idx < 0 || num < 1 {
		return
	}
	if fb.Idx == num && fb.It != nil {
		if fb.Lines.Len() == 1 {
			fb.Idx = 0
			fb.Lines.Init()
			fb.It = fb.Lines.Front()
			return
		}
		if fb.Idx == 1 {
			fb.It = fb.It.Next()
			fb.Lines.Remove(fb.It.Prev())
			return
		}
		if fb.It != fb.Lines.Back() {
			fb.It = fb.It.Next()
		} else {
			fb.It = fb.It.Prev()
			fb.Idx--
		}
		fb.Lines.Remove(a)
	}
}

func NewFileBlame() (*FileBlame, error) {
	fb := FileBlame{Idx: 0, It: &list.Element{}, Lines: list.New()}
	fb.It = fb.Lines.Front()
	fb.Idx = 0
	return &fb, nil
}
