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

// RefreshAndSaveConnection populate from request input into connection which come from REST functions to connection struct and save to DB
// and only change value which `data` has
// mergeFieldsToConnection merges fields from data
// `connection` is the pointer of a plugin connection
// `data` is http request input param
func RefreshAndSaveConnection(connection interface{}, data map[string]interface{}, db *gorm.DB) error {
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

func saveToDb(connection interface{}, db *gorm.DB) error {
	dataVal := reflect.ValueOf(connection)
	if dataVal.Kind() != reflect.Ptr {
		panic("entityPtr is not a pointer")
	}

	dataType := reflect.Indirect(dataVal).Type()
	fieldName := getEncryptField(dataType, "encrypt")
	plainPwd := ""
	err := doEncrypt(dataVal, fieldName)
	if err != nil {
		return err
	}
	err = db.Clauses(clause.OnConflict{UpdateAll: true}).Save(connection).Error
	if err != nil {
		return err
	}

	err = doDecrypt(dataVal, fieldName)
	if err != nil {
		return err
	}
	dataVal.Elem().FieldByName(fieldName).Set(reflect.ValueOf(plainPwd))

	return err
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
	// encrypt
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

// FindConnectionByInput finds connection from db  by parsing request input and decrypt it
func FindConnectionByInput(input *core.ApiResourceInput, connection interface{}, db *gorm.DB) error {
	dataVal := reflect.ValueOf(connection)
	if dataVal.Kind() != reflect.Ptr {
		return fmt.Errorf("connection is not a pointer")
	}

	id, err := GetConnectionIdByInputParam(input)
	if err != nil {
		return fmt.Errorf("invalid connectionId")
	}

	err = db.First(connection, id).Error
	if err != nil {
		fmt.Printf("--- %s", err.Error())
		return err
	}

	dataType := reflect.Indirect(dataVal).Type()

	fieldName := getEncryptField(dataType, "encrypt")
	return doDecrypt(dataVal, fieldName)

}

// GetConnectionIdByInputParam gets connectionId by parsing request input
func GetConnectionIdByInputParam(input *core.ApiResourceInput) (uint64, error) {
	connectionId := input.Params["connectionId"]
	if connectionId == "" {
		return 0, fmt.Errorf("missing connectionId")
	}
	return strconv.ParseUint(connectionId, 10, 64)
}

func getEncryptField(t reflect.Type, tag string) string {
	fieldName := ""
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			fieldName = getEncryptField(field.Type, tag)
		} else {
			if field.Tag.Get(tag) == "yes" {
				fieldName = field.Name
			}
		}
	}
	return fieldName
}

// DecryptConnection decrypts password/token field for connection
func DecryptConnection(connection interface{}, fieldName string) error {
	dataVal := reflect.ValueOf(connection)
	if dataVal.Kind() != reflect.Ptr {
		panic("connection is not a pointer")
	}
	if len(fieldName) == 0 {
		dataType := reflect.Indirect(dataVal).Type()
		fieldName = getEncryptField(dataType, "encrypt")
	}
	return doDecrypt(dataVal, fieldName)
}

func doDecrypt(dataVal reflect.Value, fieldName string) error {
	encryptCode, err := getEncKey()
	if err != nil {
		return err
	}
	if len(fieldName) > 0 {
		decryptStr, _ := core.Decrypt(encryptCode, dataVal.Elem().FieldByName(fieldName).String())

		dataVal.Elem().FieldByName(fieldName).Set(reflect.ValueOf(decryptStr))
	}
	return nil
}

func doEncrypt(dataVal reflect.Value, fieldName string) error {
	encryptCode, err := getEncKey()
	if err != nil {
		return err
	}
	if len(fieldName) > 0 {
		plainPwd := dataVal.Elem().FieldByName(fieldName).String()
		encyptedStr, err := core.Encrypt(encryptCode, plainPwd)

		if err != nil {
			return err
		}
		dataVal.Elem().FieldByName(fieldName).Set(reflect.ValueOf(encyptedStr))
	}
	return nil
}
