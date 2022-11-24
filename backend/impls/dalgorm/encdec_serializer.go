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

package dalgorm

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"

	"gorm.io/gorm/schema"
)

var _ schema.SerializerInterface = (*EncDecSerializer)(nil)

// EncDecSerializer is responsible for field encryption/decryption in Application Level
// Ref: https://gorm.io/docs/serializer.html
type EncDecSerializer struct {
	encKey string
}

// Scan implements serializer interface
func (es *EncDecSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	dbValue = models.UnwrapObject(dbValue)
	fieldValue := reflect.New(field.FieldType)
	if dbValue != nil {
		var base64str string
		switch v := dbValue.(type) {
		case []byte:
			base64str = string(v)
		case string:
			base64str = v
		default:
			return fmt.Errorf("failed to decrypt value: %#v", dbValue)
		}

		decrypted, err := plugin.Decrypt(es.encKey, base64str)
		if err != nil {
			return err
		}
		switch fieldValue.Elem().Kind() {
		case reflect.Slice:
			bytes := []byte(decrypted)
			_ = json.Unmarshal(bytes, fieldValue.Interface())
			field.ReflectValueOf(ctx, dst).Set(fieldValue.Elem())
		case reflect.String:
			field.ReflectValueOf(ctx, dst).SetString(decrypted)
		default:
			return fmt.Errorf("failed to decrypt value: %#v", dbValue)
		}
	}
	return nil
}

// Value implements serializer interface
func (es *EncDecSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	var target string
	switch v := fieldValue.(type) {
	case json.RawMessage:
		target = string(v)
	case string:
		target = v
	default:
		return nil, fmt.Errorf("failed to encrypt value: %#v", fieldValue)
	}
	return plugin.Encrypt(es.encKey, target)
}

// Init the encdec serializer
func Init(encKey string) {
	schema.RegisterSerializer("encdec", &EncDecSerializer{encKey: encKey})
}
