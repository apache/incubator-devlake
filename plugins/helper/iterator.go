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
	"database/sql"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core/dal"
)

// Iterator FIXME ...
type Iterator interface {
	HasNext() bool
	Fetch() (interface{}, error)
	Close() error
}

// DalCursorIterator FIXME ...
type DalCursorIterator struct {
	db       dal.Dal
	cursor   *sql.Rows
	elemType reflect.Type
}

// NewDalCursorIterator FIXME ...
func NewDalCursorIterator(db dal.Dal, cursor *sql.Rows, elemType reflect.Type) (*DalCursorIterator, error) {
	return &DalCursorIterator{
		db:       db,
		cursor:   cursor,
		elemType: elemType,
	}, nil
}

// HasNext FIXME ...
func (c *DalCursorIterator) HasNext() bool {
	return c.cursor.Next()
}

// Fetch FIXME ...
func (c *DalCursorIterator) Fetch() (interface{}, error) {
	elem := reflect.New(c.elemType).Interface()
	err := c.db.Fetch(c.cursor, elem)
	if err != nil {
		return nil, err
	}
	return elem, nil
}

// Close interator
func (c *DalCursorIterator) Close() error {
	return c.cursor.Close()
}

var _ Iterator = (*DalCursorIterator)(nil)

// DateIterator FIXME ...
type DateIterator struct {
	startTime time.Time
	endTime   time.Time
	Days      int
	Current   int
}

// DatePair FIXME ...
type DatePair struct {
	PairStartTime time.Time
	PairEndTime   time.Time
}

// HasNext FIXME ...
func (c *DateIterator) HasNext() bool {
	return c.Current < c.Days
}

// Fetch FIXME ...
func (c *DateIterator) Fetch() (interface{}, error) {
	c.Current++
	return &DatePair{
		PairStartTime: c.startTime.AddDate(0, 0, c.Current),
		PairEndTime:   c.endTime.AddDate(0, 0, c.Current),
	}, nil
}

// Close iterator
func (c *DateIterator) Close() error {
	return nil
}

// NewDateIterator FIXME ...
func NewDateIterator(days int) (*DateIterator, error) {
	endTime := time.Now().Truncate(24 * time.Hour)
	return &DateIterator{
		startTime: endTime.AddDate(0, 0, -days-1),
		endTime:   endTime.AddDate(0, 0, -days),
		Days:      days,
		Current:   0,
	}, nil
}

type QueueIteratorNode struct {
	data interface{}
	next *QueueIteratorNode
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

type QueueIterator struct {
	queue *models.Queue
}

func (q *QueueIterator) HasNext() bool {
	return q.queue.GetCount() > 0
}

func (q *QueueIterator) Fetch() (interface{}, error) {
	return q.queue.PullWithOutLock(), nil
}

func (q *QueueIterator) Push(data models.QueueNode) {
	q.queue.PushWitouLock(data)
}

func (q *QueueIterator) Close() error {
	q.queue.CleanWithOutLock()
	return nil
}

func NewQueueIterator() *QueueIterator {
	return &QueueIterator{
		queue: models.NewQueue(),
	}
}
