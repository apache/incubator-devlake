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
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"reflect"
	"strings"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
)

// BatchSave performs mulitple records persistence of a specific type in one sql query to improve the performance
type BatchSave struct {
	basicRes core.BasicRes
	log      core.Logger
	db       dal.Dal
	slotType reflect.Type
	// slots can not be []interface{}, because gorm wouldn't take it
	// I'm guessing the reason is the type information lost when converted to interface{}
	slots      reflect.Value
	current    int
	size       int
	valueIndex map[string]int
	primaryKey []reflect.StructField
}

// NewBatchSave creates a new BatchSave instance
func NewBatchSave(basicRes core.BasicRes, slotType reflect.Type, size int) (*BatchSave, errors.Error) {
	if slotType.Kind() != reflect.Ptr {
		return nil, errors.Default.New("slotType must be a pointer")
	}
	db := basicRes.GetDal()
	primaryKey := db.GetPrimaryKeyFields(slotType)
	// check if it have primaryKey
	if len(primaryKey) == 0 {
		return nil, errors.Default.New(fmt.Sprintf("%s no primary key", slotType.String()))
	}

	log := basicRes.GetLogger().Nested(slotType.String())
	return &BatchSave{
		basicRes:   basicRes,
		log:        log,
		db:         db,
		slotType:   slotType,
		slots:      reflect.MakeSlice(reflect.SliceOf(slotType), size, size),
		size:       size,
		valueIndex: make(map[string]int),
		primaryKey: primaryKey,
	}, nil
}

// Add record to cache. BatchSave would flush them into Database when cache is max out
func (c *BatchSave) Add(slot interface{}) errors.Error {
	// type checking
	if reflect.TypeOf(slot) != c.slotType {
		return errors.Default.New("sub cache type mismatched")
	}
	if reflect.ValueOf(slot).Kind() != reflect.Ptr {
		return errors.Default.New("slot is not a pointer")
	}
	// deduplication
	key := getKeyValue(slot, c.primaryKey)

	if key != "" {
		if index, ok := c.valueIndex[key]; !ok {
			c.valueIndex[key] = c.current
		} else {
			c.slots.Index(index).Set(reflect.ValueOf(slot))
			return nil
		}
	}
	c.slots.Index(c.current).Set(reflect.ValueOf(slot))
	c.current++
	// flush out into database if max outed
	if c.current == c.size {
		return c.Flush()
	} else if c.current%100 == 0 {
		c.log.Debug("batch save current: %d", c.current)
	}
	return nil
}

// Flush save cached records into database
func (c *BatchSave) Flush() errors.Error {
	err := c.db.CreateOrUpdate(c.slots.Slice(0, c.current).Interface())
	if err != nil {
		return err
	}
	c.log.Debug("batch save flush total %d records to database", c.current)
	c.current = 0
	c.valueIndex = make(map[string]int)
	return nil
}

// Close would flash the cache and release resources
func (c *BatchSave) Close() errors.Error {
	if c.current > 0 {
		return c.Flush()
	}
	return nil
}

func getKeyValue(iface interface{}, primaryKey []reflect.StructField) string {
	var ss []string
	ifv := reflect.ValueOf(iface)
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
	}
	for _, key := range primaryKey {
		v := ifv.FieldByName(key.Name)
		s := fmt.Sprintf("%v", v.Interface())
		if s != "" {
			ss = append(ss, s)
		}
	}
	return strings.Join(ss, ":")
}
