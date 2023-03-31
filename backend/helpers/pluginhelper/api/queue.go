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

package api

import (
	"sync"
	"time"
)

// QueueNode represents a node in the queue
type QueueNode interface {
	Next() interface{}
	SetNext(next interface{})
	Data() interface{}
}

// Queue represetns a queue
type Queue struct {
	count   int64
	head    QueueNode
	tail    QueueNode
	mux     sync.Mutex
	working int64 // working count
}

// reduce working count
func (q *Queue) Finish(count int64) {
	q.mux.Lock()
	defer q.mux.Unlock()

	q.working -= count
}

// Push add a node to queue
func (q *Queue) Push(node QueueNode) {
	q.mux.Lock()
	defer q.mux.Unlock()

	q.PushWithoutLock(node)
}

// Pull get a node from queue
// it will add the working count and blocked when there are no node on queue but working count not zero
func (q *Queue) Pull() QueueNode {
	q.mux.Lock()
	defer q.mux.Unlock()
	node := q.PullWithOutLock()
	if node != nil {
		return node
	} else {
		return nil
	}
}

func (q *Queue) PullWithWorkingBlock() QueueNode {
	q.mux.Lock()
	defer q.mux.Unlock()

	for {
		node := q.PullWithOutLock()
		if node != nil {
			q.working++

			return node
		} else if q.working > 0 {
			q.mux.Unlock()

			time.Sleep(time.Second)

			q.mux.Lock()
		} else {
			return nil
		}
	}
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
	var node QueueNode

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

// GetCount get the node count in query and only return zero when working zero
func (q *Queue) GetCountWithWorkingBlock() int64 {
	q.mux.Lock()
	defer q.mux.Unlock()

	for {
		if q.count == 0 && q.working > 0 {
			q.mux.Unlock()

			time.Sleep(time.Second)

			q.mux.Lock()
		} else {
			return q.count
		}
	}
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

// QueueIteratorNode implements the api.Iterator interface with ability to accept new item when being iterated
type QueueIteratorNode struct {
	next QueueNode
	data interface{}
}

// Next return the next node
func (q *QueueIteratorNode) Next() interface{} {
	if q.next == nil {
		return nil
	}
	return q.next
}

// SetNext updates the next pointer of the node
func (q *QueueIteratorNode) SetNext(next interface{}) {
	if next == nil {
		q.next = nil
	} else {
		q.next, _ = next.(QueueNode)
	}
}

// Data returns data of the node
func (q *QueueIteratorNode) Data() interface{} {
	return q.data
}

// NewQueueIteratorNode creates a new QueueIteratorNode
func NewQueueIteratorNode(data interface{}) *QueueIteratorNode {
	return &QueueIteratorNode{
		data: data,
	}
}
