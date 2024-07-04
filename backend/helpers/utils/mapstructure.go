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
	"encoding/json"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/spf13/cast"

	"github.com/apache/incubator-devlake/core/errors"

	"github.com/mitchellh/mapstructure"
)

func decodeHookStringFloat64(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() == reflect.Float64 && t.Kind() == reflect.Struct && t == reflect.TypeOf(common.StringFloat64{}) {
		return *common.NewStringFloat64FromAny(data), nil
	}
	if f.Kind() == reflect.Float64 && t.Kind() == reflect.Ptr && t == reflect.TypeOf(&common.StringFloat64{}) {
		return common.NewStringFloat64FromAny(data), nil
	}
	return data, nil
}

func decodeHookStringToTime(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() == reflect.String {
		if t.Kind() == reflect.Struct && t == reflect.TypeOf(time.Time{}) {
			timeData, err := common.ConvertStringToTime(cast.ToString(data))
			if err != nil {
				return data, err
			}
			return timeData, nil
		}
		if t.Kind() == reflect.Ptr && t == reflect.TypeOf(&time.Time{}) {
			timeData, err := common.ConvertStringToTime(cast.ToString(data))
			if err != nil {
				return data, err
			}
			return &timeData, nil
		}
	}
	return data, nil
}

func zeroesSlice(from reflect.Value, to reflect.Value) (interface{}, error) {
	if from.Kind() == reflect.Slice && from.IsNil() {
		return from.Interface(), nil
	}

	if to.CanSet() && to.Kind() == reflect.Slice {
		to.Set(reflect.MakeSlice(reflect.SliceOf(to.Type().Elem()), 0, 0))
	}

	return from.Interface(), nil
}

func decodeUrlValues(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	// to support decoding url.Values (query string) to non-array variables
	if to.Kind() != reflect.Slice && to.Kind() != reflect.Array &&
		(from.Kind() == reflect.Slice || from.Kind() == reflect.Array) {
		v := reflect.ValueOf(data)
		if v.Len() == 1 {
			data = v.Index(0).Interface()
			var result interface{}
			err := DecodeMapStruct(data, &result, true)
			return result, err
		}
	}
	return data, nil
}

func decodeJsonRawMessage(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if to == reflect.TypeOf(json.RawMessage{}) {
		return json.Marshal(data)
	}
	return data, nil
}

func decodeIso8601Time(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if data == nil {
		return nil, nil
	}
	if to != reflect.TypeOf(common.Iso8601Time{}) && to != reflect.TypeOf(time.Time{}) {
		return data, nil
	}

	var tt time.Time
	var err error
	switch from.Kind() {
	case reflect.String:
		tt, err = common.ConvertStringToTime(data.(string))
	case reflect.Float64:
		tt = time.Unix(0, int64(data.(float64))*int64(time.Millisecond))
	case reflect.Int64:
		tt = time.Unix(0, data.(int64)*int64(time.Millisecond))
	case reflect.Struct:
		if from == reflect.TypeOf(time.Time{}) || from == reflect.TypeOf(common.Iso8601Time{}) {
			return data, nil
		}
	case reflect.Ptr:
		if from == reflect.TypeOf(&time.Time{}) || from == reflect.TypeOf(&common.Iso8601Time{}) {
			return data, nil
		}
	}
	if err != nil {
		return data, nil
	}
	switch to {
	case reflect.TypeOf(common.Iso8601Time{}):
		return common.Iso8601Time{Time: tt}, nil
	case reflect.TypeOf(&common.Iso8601Time{}):
		return &common.Iso8601Time{Time: tt}, nil
	}
	return data, nil
}

// DecodeMapStruct with time.Time and Iso8601Time support
func DecodeMapStruct(input interface{}, result interface{}, zeroFields bool) errors.Error {
	result = models.UnwrapObject(result)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ZeroFields: zeroFields,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			decodeHookStringFloat64,
			decodeHookStringToTime,
			decodeUrlValues,
			decodeJsonRawMessage,
			decodeIso8601Time,
			zeroesSlice,
		),
		Result:           result,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return errors.Convert(err)
	}

	if err := decoder.Decode(input); err != nil {
		return errors.Convert(err)
	}
	return errors.Convert(err)
}
