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
	"encoding/base64"
	"fmt"
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"strconv"
)

type BaseConnection struct {
	Name string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	common.Model
}

type BasicAuth struct {
	Username string `mapstructure:"username" validate:"required" json:"username"`
	Password string `mapstructure:"password" validate:"required" json:"password" encrypt:"yes"`
}

func (ba BasicAuth) GetEncodedToken() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", ba.Username, ba.Password)))
}

type AccessToken struct {
	Token string `mapstructure:"token" validate:"required" json:"token" encrypt:"yes"`
}

type RestConnection struct {
	BaseConnection `mapstructure:",squash"`
	Endpoint       string `mapstructure:"endpoint" validate:"required" json:"endpoint"`
	Proxy          string `mapstructure:"proxy" json:"proxy"`
	RateLimit      int    `comment:"api request rate limt per hour" json:"rateLimit"`
}

// CreateConnection populate from request input into connection which come from REST functions to connection struct and save to DB
// and only change value which `data` has
// mergeFieldsToConnection merges fields from data
// `connection` is the pointer of a plugin connection
// `data` is http request input param
func CreateConnection(data map[string]interface{}, connection interface{}, db *gorm.DB) error {
	var err error
	// update fields from request body
	err = mergeFieldsToConnection(connection, data)
	if err != nil {
		return err
	}
	err = saveToDb(connection, db)
	if err != nil {
		return err
	}
	return nil
}

func PatchConnection(input *core.ApiResourceInput, connection interface{}, db *gorm.DB) error {
	err := GetConnection(input.Params, connection, db)
	if err != nil {
		return err
	}

	err = CreateConnection(input.Body, connection, db)
	if err != nil {
		return err
	}

	return nil
}

func saveToDb(connection interface{}, db *gorm.DB) error {
	err := EncryptConnection(connection)
	if err != nil {
		return err
	}
	err = db.Clauses(clause.OnConflict{UpdateAll: true}).Save(connection).Error
	if err != nil {
		return err
	}

	return DecryptConnection(connection)
}

// mergeFieldsToConnection will populate all value in map to connection struct and validate the struct
func mergeFieldsToConnection(specificConnection interface{}, connections ...map[string]interface{}) error {
	// decode
	for _, connection := range connections {
		err := mapstructure.Decode(connection, specificConnection)
		if err != nil {
			return err
		}
	}
	// validate
	vld := validator.New()
	err := vld.Struct(specificConnection)
	if err != nil {
		return err
	}

	return nil
}

func getEncKey() (string, error) {
	v := config.GetConfig()
	encKey := v.GetString(core.EncodeKeyEnvStr)
	if encKey == "" {
		// Randomly generate a bunch of encryption keys and set them to config
		encKey = core.RandomEncKey()
		v.Set(core.EncodeKeyEnvStr, encKey)
		err := config.WriteConfig(v)
		if err != nil {
			return encKey, err
		}
	}
	return encKey, nil
}

// GetConnection finds connection from db  by parsing request input and decrypt it
func GetConnection(data map[string]string, connection interface{}, db *gorm.DB) error {
	id, err := GetConnectionIdByInputParam(data)
	if err != nil {
		return fmt.Errorf("invalid connectionId")
	}

	err = db.First(connection, id).Error
	if err != nil {
		return err
	}

	return DecryptConnection(connection)

}

// ListConnections returns all connections with password/token decrypted
func ListConnections(connections interface{}, db *gorm.DB) error {
	err := db.Find(connections).Error
	connPtr := reflect.ValueOf(connections)
	connVal := reflect.Indirect(connPtr)
	if err != nil {
		return err
	}
	for i := 0; i < connVal.Len(); i++ {
		//connVal.Index(i) returns value of ith elem in connections, .Elem() reutrns the original elem
		tmp := connVal.Index(i).Elem()
		err = DecryptConnection(tmp.Addr().Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

// GetConnectionIdByInputParam gets connectionId by parsing request input
func GetConnectionIdByInputParam(data map[string]string) (uint64, error) {
	connectionId := data["connectionId"]
	if connectionId == "" {
		return 0, fmt.Errorf("missing connectionId")
	}
	return strconv.ParseUint(connectionId, 10, 64)
}

func firstFieldNameWithTag(t reflect.Type, tag string) string {
	fieldName := ""
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			fieldName = firstFieldNameWithTag(field.Type, tag)
		} else {
			if field.Tag.Get(tag) == "yes" {
				fieldName = field.Name
			}
		}
	}
	return fieldName
}

// DecryptConnection decrypts password/token field for connection
func DecryptConnection(connection interface{}) error {
	dataVal := reflect.ValueOf(connection)
	if dataVal.Kind() != reflect.Ptr {
		panic("connection is not a pointer")
	}
	encKey, err := getEncKey()
	if err != nil {
		return nil
	}
	dataType := reflect.Indirect(dataVal).Type()
	fieldName := firstFieldNameWithTag(dataType, "encrypt")
	if len(fieldName) > 0 {
		decryptStr, _ := core.Decrypt(encKey, dataVal.Elem().FieldByName(fieldName).String())
		dataVal.Elem().FieldByName(fieldName).Set(reflect.ValueOf(decryptStr))
	}
	return nil
}

func EncryptConnection(connection interface{}) error {
	dataVal := reflect.ValueOf(connection)
	if dataVal.Kind() != reflect.Ptr {
		panic("connection is not a pointer")
	}
	encKey, err := getEncKey()
	if err != nil {
		return err
	}
	dataType := reflect.Indirect(dataVal).Type()
	fieldName := firstFieldNameWithTag(dataType, "encrypt")
	if len(fieldName) > 0 {
		plainPwd := dataVal.Elem().FieldByName(fieldName).String()
		encyptedStr, err := core.Encrypt(encKey, plainPwd)

		if err != nil {
			return err
		}
		dataVal.Elem().FieldByName(fieldName).Set(reflect.ValueOf(encyptedStr))
	}
	return nil
}
