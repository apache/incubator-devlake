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

package utils

import (
	"fmt"
	"regexp"
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

// GetTimeFieldFromMap retrieves a time field from a map.
// allFields: A map containing all fields.
// fieldName: The name of the field to retrieve.
// loc: The timezone location.
// Returns:
//
//	*time.Time: A pointer to the time.Time if the field exists and can be converted to time.Time.
//	error: An error if the field does not exist or an error occurs.
func GetTimeFieldFromMap(allFields map[string]interface{}, fieldName string, loc *time.Location) (*time.Time, error) {
	val, ok := allFields[fieldName]
	if !ok {
		return nil, fmt.Errorf("Field %s not found", fieldName)
	}
	var temp time.Time
	switch v := val.(type) {
	case string:
		if v == "" || v == "null" || v == "{}" {
			// In Zentao, the field `deadline`'s value may be "{}".
			return nil, nil
		}
		// If value is a string with the format yyyy-MM-dd, use loc to parse it
		isDateFormat, _ := regexp.Match(`\d{4}-\d{2}-\d{2}`, []byte(v))
		if isDateFormat {
			temp, _ = common.ConvertStringToTimeInLoc(v, loc)
		} else {
			temp, _ = common.ConvertStringToTime(v)
		}

	default:
		temp, _ = v.(time.Time)
	}
	if temp.IsZero() {
		return nil, nil
	}
	return &temp, nil
}
