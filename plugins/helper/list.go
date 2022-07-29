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

package helper

// ListBaseNode 'abstract' base struct for Nodes that are chained in a linked list manner
type ListBaseNode struct {
	next QueueNode
}

func (l *ListBaseNode) Data() interface{} {
	// default implementation
	return nil
}

func (l *ListBaseNode) Next() interface{} {
	if l.next == nil {
		return nil
	}
	return l.next
}

func (l *ListBaseNode) SetNext(next interface{}) {
	if next == nil {
		l.next = nil
	} else {
		l.next = next.(QueueNode)
	}
}

// NewListBaseNode create and init a new node (only to be called by subclasses)
func NewListBaseNode() *ListBaseNode {
	return &ListBaseNode{
		next: nil,
	}
}

// check if is all right for interface QueueNode
var _ QueueNode = (*ListBaseNode)(nil)
