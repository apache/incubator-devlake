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

import (
	"sync"
	"sync/atomic"
)

type QueueNode interface {
	Next() interface{}
	SetNext(next interface{})
	Data() interface{}
}

type Queue struct {
	count int64
	head  QueueNode
	tail  QueueNode
	mux   sync.Mutex
}

// Push add a node to queue
func (q *Queue) Push(node QueueNode) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.PushWithoutLock(node)
}

// Pull get a node from queue
func (q *Queue) Pull(add *int64) QueueNode {
	q.mux.Lock()
	defer q.mux.Unlock()

	node := q.PullWithOutLock()

	if node == nil {
		return nil
	}
	if add != nil {
		atomic.AddInt64(add, 1)
	}
	return node
}

// PushWithoutLock is no lock mode of Push
func (q *Queue) PushWithoutLock(node QueueNode) {
	if q.tail == nil {
		q.head = node
		q.tail = node
		q.count = 1
	} else {
		q.tail.SetNext(node)
		q.tail = node
		q.count++
	}
}

// PullWithOutLock is no lock mode of Pull
func (q *Queue) PullWithOutLock() QueueNode {
	var node QueueNode = nil

	if q.head != nil {
		node = q.head
		q.head, _ = node.Next().(QueueNode)

		if q.head == nil {
			q.tail = nil
		}

		node.SetNext(nil)
		q.count--
	} else {
		q.count = 0
	}
	return node
}

// CleanWithOutLock is no lock mode of Clean
func (q *Queue) CleanWithOutLock() {
	q.count = 0
	q.head = nil
	q.tail = nil
}

// Clean remove all node on queue
func (q *Queue) Clean() {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.CleanWithOutLock()
}

// GetCountWithOutLock is no lock mode of GetCount
func (q *Queue) GetCountWithOutLock() int64 {
	return q.count
}

// GetCount get the node count
func (q *Queue) GetCount() int64 {
	q.mux.Lock()
	defer q.mux.Unlock()
	return q.count
}

// NewQueue create and init a new Queue
func NewQueue() *Queue {
	return &Queue{
		count: int64(0),
		head:  nil,
		tail:  nil,
		mux:   sync.Mutex{},
	}
}

type QueueIteratorNode struct {
	next *QueueIteratorNode
	data interface{}
}

func (q *QueueIteratorNode) Next() interface{} {
	if q.next == nil {
		return nil
	}
	return q.next
}

func (q *QueueIteratorNode) SetNext(next interface{}) {
	q.next, _ = next.(*QueueIteratorNode)
}

func (q *QueueIteratorNode) Data() interface{} {
	return q.data
}

func NewQueueIteratorNode(data interface{}) *QueueIteratorNode {
	return &QueueIteratorNode{
		data: data,
	}
}
