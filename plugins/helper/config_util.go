package helper

import (
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"reflect"
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
		panic("input is not a pointer")
	}
	tf := reflect.Indirect(vf).Type()
	fieldTags := make([]string, 0)
	fieldNames := make([]string, 0)
	fieldTypes := make([]reflect.Type, 0)
	walkFields(tf, &fieldNames, &fieldTypes, &fieldTags, tag)
	length := len(fieldNames)
	for i := 0; i < length; i++ {
		fieldName := fieldNames[i]
		fieldType := fieldTypes[i]
		fieldTag := fieldTags[i]

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
		switch fieldType.String() {
		case "string":
			output.Set(fieldTag, vfField.String())
		case "int", "int64":
			output.Set(fieldTag, vfField.Int())
		case "float64":
			output.Set(fieldTag, vfField.Float())
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
		panic("output is not a pointer")
	}
	tf := reflect.Indirect(vf).Type()
	fieldTags := make([]string, 0)
	fieldNames := make([]string, 0)
	fieldTypes := make([]reflect.Type, 0)
	walkFields(tf, &fieldNames, &fieldTypes, &fieldTags, tag)
	length := len(fieldNames)
	for i := 0; i < length; i++ {
		fieldName := fieldNames[i]
		fieldType := fieldTypes[i]
		fieldTag := fieldTags[i]

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
		switch fieldType.String() {
		case "string":
			vfField.SetString(input.GetString(fieldTag))
		case "int", "int64":
			vfField.SetInt(input.GetInt64(fieldTag))
		case "float64":
			vfField.SetFloat(input.GetFloat64(fieldTag))
		default:
		}
	}
	return nil
}

func walkFields(t reflect.Type, fieldNames *[]string, fieldTypes *[]reflect.Type, fieldTags *[]string, tag string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			walkFields(field.Type, fieldNames, fieldTypes, fieldTags, tag)
		} else {
			fieldTag := field.Tag.Get(tag)
			*fieldNames = append(*fieldNames, field.Name)
			*fieldTypes = append(*fieldTypes, field.Type)
			*fieldTags = append(*fieldTags, fieldTag)
		}
	}
}
