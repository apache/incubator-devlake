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
}

const BATCH_SAVE_UPDATE_ONLY = 0

// NewBatchSave creates a new BatchSave instance
func NewBatchSave(basicRes core.BasicRes, slotType reflect.Type, size int) (*BatchSave, error) {
	if slotType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("slotType must be a pointer")
	}
	if !hasPrimaryKey(slotType) {
		return nil, fmt.Errorf("%s no primary key", slotType.String())
	}
	log := basicRes.GetLogger().Nested(slotType.String())
	return &BatchSave{
		basicRes:   basicRes,
		log:        log,
		db:         basicRes.GetDal(),
		slotType:   slotType,
		slots:      reflect.MakeSlice(reflect.SliceOf(slotType), size, size),
		size:       size,
		valueIndex: make(map[string]int),
	}, nil
}

// Add record to cache. BatchSave would flush them into Database when cache is max out
func (c *BatchSave) Add(slot interface{}) error {
	// type checking
	if reflect.TypeOf(slot) != c.slotType {
		return fmt.Errorf("sub cache type mismatched")
	}
	if reflect.ValueOf(slot).Kind() != reflect.Ptr {
		return fmt.Errorf("slot is not a pointer")
	}
	// deduplication
	key := getPrimaryKeyValue(slot)

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
func (c *BatchSave) Flush() error {
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
func (c *BatchSave) Close() error {
	if c.current > 0 {
		return c.Flush()
	}
	return nil
}

func isPrimaryKey(f reflect.StructField) bool {
	tag := strings.TrimSpace(f.Tag.Get("gorm"))
	return strings.HasPrefix(strings.ToLower(tag), "primarykey")
}

func hasPrimaryKey(ifv reflect.Type) bool {
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
	}
	for i := 0; i < ifv.NumField(); i++ {
		v := ifv.Field(i)
		if ok := isPrimaryKey(v); ok {
			return true
		} else if v.Type.Kind() == reflect.Struct {
			if ok := hasPrimaryKey(v.Type); ok {
				return true
			}
		}
	}
	return false
}

func getPrimaryKeyValue(iface interface{}) string {
	var ss []string
	ifv := reflect.ValueOf(iface)
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
	}
	for i := 0; i < ifv.NumField(); i++ {
		v := ifv.Field(i)
		if isPrimaryKey(ifv.Type().Field(i)) {
			s := fmt.Sprintf("%v", v.Interface())
			if s != "" {
				ss = append(ss, s)
			}
		} else if v.Kind() == reflect.Struct {
			s := getPrimaryKeyValue(v.Interface())
			if s != "" {
				ss = append(ss, s)
			}
		}
	}
	return strings.Join(ss, ":")
}
