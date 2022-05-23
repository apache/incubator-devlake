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
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type CSTTime time.Time

func (jt *CSTTime) UnmarshalJSON(b []byte) error {
	timeString := string(b)
	if timeString == "null" {
		return nil
	}
	if strings.Contains(timeString, "0000-00-00") {
		return nil
	}
	timeString = strings.Trim(timeString, `"`)
	if len(timeString) == 10 {
		timeString = fmt.Sprintf("%s 00:00:00", timeString)
	}
	cstZone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeString, cstZone)
	if err != nil {
		return err
	}
	*jt = CSTTime(t)
	return nil
}

func (jt CSTTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	t := (time.Time)(jt)
	if t.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t, nil
}
func (jt *CSTTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*jt = CSTTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
