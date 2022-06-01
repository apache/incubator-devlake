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
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/plugins/core/dal"
)

type Iterator interface {
	HasNext() bool
	Fetch() (interface{}, error)
	Close() error
}

type CursorIterator struct {
	db       dal.Dal
	cursor   dal.Cursor
	elemType reflect.Type
}

func NewCursorIterator(db dal.Dal, cursor dal.Cursor, elemType reflect.Type) (*CursorIterator, error) {
	return &CursorIterator{
		db:       db,
		cursor:   cursor,
		elemType: elemType,
	}, nil
}

func (c *CursorIterator) HasNext() bool {
	return c.cursor.Next()
}

func (c *CursorIterator) Fetch() (interface{}, error) {
	elem := reflect.New(c.elemType).Interface()
	err := c.db.Fetch(c.cursor, elem)
	if err != nil {
		return nil, err
	}
	return elem, nil
}

func (c *CursorIterator) Close() error {
	return c.cursor.Close()
}

var _ Iterator = (*CursorIterator)(nil)

type DateIterator struct {
	startTime time.Time
	endTime   time.Time
	Days      int
	Current   int
}

type DatePair struct {
	PairStartTime time.Time
	PairEndTime   time.Time
}

func (c *DateIterator) HasNext() bool {
	return c.Current < c.Days
}

func (c *DateIterator) Fetch() (interface{}, error) {
	c.Current++
	return &DatePair{
		PairStartTime: c.startTime.AddDate(0, 0, c.Current),
		PairEndTime:   c.endTime.AddDate(0, 0, c.Current),
	}, nil

}

func (c *DateIterator) Close() error {
	return nil
}

func NewDateIterator(days int) (*DateIterator, error) {
	endTime := time.Now().Truncate(24 * time.Hour)
	return &DateIterator{
		startTime: endTime.AddDate(0, 0, -days-1),
		endTime:   endTime.AddDate(0, 0, -days),
		Days:      days,
		Current:   0,
	}, nil
}
