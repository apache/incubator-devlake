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

	"github.com/apache/incubator-devlake/utils"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// DecodeStruct validates `input` struct with `validator` and set it into viper
// `tag` represent the fields when setting config, and the fields with `tag` shall prevail.
// `input` must be a pointer
func DecodeStruct(output *viper.Viper, input interface{}, data map[string]interface{}, tag string) error {
	// update fields from request body
	err := mapstructure.Decode(data, input)
	if err != nil {
		return err
	}
	// validate fields with `validate` tag
	vld := validator.New()
	err = vld.Struct(input)
	if err != nil {
		return err
	}
	err = validator.New().Struct(input)
	if err != nil {
		return err
	}

	vf := reflect.ValueOf(input)
	if vf.Kind() != reflect.Ptr {
		return fmt.Errorf("input %v is not a pointer", input)
	}

	for _, f := range utils.WalkFields(reflect.Indirect(vf).Type(), nil) {
		fieldName := f.Name
		fieldType := f.Type
		fieldTag := f.Tag.Get(tag)

		// Check if the first letter is uppercase (indicates a public element, accessible)
		ascii := rune(fieldName[0])
		if int(ascii) < int('A') || int(ascii) > int('Z') {
			continue
		}

		// View their tags in order to filter out members who don't have a valid tag set
		if fieldTag == "" {
			continue
		}
		vfField := vf.Elem().FieldByName(fieldName)
		switch fieldType.Kind() {
		case reflect.String:
			output.Set(fieldTag, vfField.String())
		case reflect.Int, reflect.Int64:
			output.Set(fieldTag, vfField.Int())
		case reflect.Float64, reflect.Float32:
			output.Set(fieldTag, vfField.Float())
		case reflect.Bool:
			output.Set(fieldTag, vfField.Bool())
		case reflect.Slice:
			elemType := vfField.Type().Elem().Kind()
			switch elemType {
			case reflect.String:
				output.Set(fieldTag, vfField.Interface().([]string))
			case reflect.Int:
				output.Set(fieldTag, vfField.Interface().([]int))
			}
		case reflect.Map:
			keyType := vfField.Type().Key().Kind()
			elemType := vfField.Type().Elem().Kind()
			if keyType == reflect.String {
				switch elemType {
				case reflect.String:
					output.Set(fieldTag, vfField.Interface().(map[string]string))
				case reflect.Interface:
					output.Set(fieldTag, vfField.Interface().(map[string]interface{}))
				}
			}
		default:
		}
	}
	return nil
}

// EncodeStruct encodes struct from viper
// `tag` represent the fields when setting config, and the fields with `tag` shall prevail.
// `object` must be a pointer
func EncodeStruct(input *viper.Viper, output interface{}, tag string) error {
	vf := reflect.ValueOf(output)
	if vf.Kind() != reflect.Ptr {
		return fmt.Errorf("output %v is not a pointer", output)
	}

	for _, f := range utils.WalkFields(reflect.Indirect(vf).Type(), nil) {
		fieldName := f.Name
		fieldType := f.Type
		fieldTag := f.Tag.Get(tag)

		// Check if the first letter is uppercase (indicates a public element, accessible)
		ascii := rune(fieldName[0])
		if int(ascii) < int('A') || int(ascii) > int('Z') {
			continue
		}

		// View their tags in order to filter out members who don't have a valid tag set
		if fieldTag == "" {
			continue
		}
		vfField := vf.Elem().FieldByName(fieldName)
		switch fieldType.Kind() {
		case reflect.String:
			vfField.SetString(input.GetString(fieldTag))
		case reflect.Int, reflect.Int64:
			vfField.SetInt(input.GetInt64(fieldTag))
		case reflect.Float64:
			vfField.SetFloat(input.GetFloat64(fieldTag))
		case reflect.Bool:
			vfField.SetBool(input.GetBool(fieldTag))
		case reflect.Slice:
			elem := vfField.Type().Elem()
			switch elem.Kind() {
			case reflect.String:
				value := input.GetStringSlice(fieldTag)
				stringSlice := reflect.MakeSlice(reflect.SliceOf(elem), 0, len(value))
				for _, item := range value {
					stringSlice = reflect.Append(stringSlice, reflect.ValueOf(item))
				}
				vfField.Set(stringSlice)
			case reflect.Int:
				value := input.GetIntSlice(fieldTag)
				intSlice := reflect.MakeSlice(reflect.SliceOf(elem), 0, len(value))
				for _, item := range value {
					intSlice = reflect.Append(intSlice, reflect.ValueOf(item))
				}
				vfField.Set(intSlice)
			}
		case reflect.Map:
			key := vfField.Type().Key()
			elem := vfField.Type().Elem()
			if key.Kind() == reflect.String {
				mapType := reflect.MapOf(key, elem)
				data := reflect.MakeMap(mapType)
				switch elem.Kind() {
				case reflect.String:
					for k, value := range input.GetStringMapString(fieldTag) {
						data.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(value))
					}
				case reflect.Interface:
					for k, value := range input.GetStringMap(fieldTag) {
						data.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(value))
					}
				}
				vfField.Set(data)
			}
		default:
		}
	}
	return nil
}
