package helper

import (
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"reflect"
)

// SaveToConifgWithMap will do two things, decode map to a struct and then save struct to config
// `object` must be a pointer
func SaveToConifgWithMap(v *viper.Viper, object interface{}, data map[string]interface{}, Tags ...string) error {
	// update fields from request body
	err := decodeMap(object, data)
	if err != nil {
		return err
	}
	err = SaveToConfig(v, object, Tags...)
	if err != nil {
		return err
	}
	return nil
}

func decodeMap(object interface{}, data map[string]interface{}) error {
	// decode
	err := mapstructure.Decode(data, object)
	if err != nil {
		return err
	}
	// validate
	vld := validator.New()
	err = vld.Struct(object)
	if err != nil {
		return err
	}

	return nil
}

// SaveToConfig Set a struct for viper
// `Tags` represent the fields when setting config, and the fields in Tags shall prevail. `Tags` that appear first have higher priority.
func SaveToConfig(v *viper.Viper, object interface{}, Tags ...string) error {
	err := validator.New().Struct(object)
	if err != nil {
		return err
	}

	vf := reflect.ValueOf(object)
	tf := reflect.Indirect(vf).Type()
	for i := 0; i < tf.NumField(); i++ {
		if vf.Elem().Field(i).IsZero() {
			continue
		}
		tfield := tf.Field(i)

		// Check if the first letter is uppercase (indicates a public element, accessible)
		ascii := rune(tfield.Name[0])
		if int(ascii) < int('A') || int(ascii) > int('Z') {
			continue
		}

		// View their tags in order to filter out members who don't have a valid tag set
		ft := ""
		for _, Tag := range Tags {
			ft = tfield.Tag.Get(Tag)
			if ft != "" {
				break
			}
		}
		if ft == "" {
			continue
		}
		vfField := vf.Elem().Field(i)
		switch vfField.Type().String() {
		case "string":
			v.Set(ft, vfField.String())
		case "int", "int64":
			v.Set(ft, vfField.Int())
		case "float64":
			v.Set(ft, vfField.Float())
		default:
		}
	}
	return v.WriteConfig()
}

// LoadFromConfig Load struct from viper
// `Tags` represent the fields when setting config, and the fields in Tags shall prevail. `Tags` that appear first have higher priority.
// `object` must be a pointer
func LoadFromConfig(v *viper.Viper, object interface{}, Tags ...string) (interface{}, error) {
	vf := reflect.ValueOf(object)
	tf := reflect.Indirect(vf).Type()
	for i := 0; i < tf.NumField(); i++ {
		tfield := tf.Field(i)

		// Check if the first letter is uppercase (indicates a public element, accessible)
		ascii := rune(tfield.Name[0])
		if int(ascii) < int('A') || int(ascii) > int('Z') {
			continue
		}

		// View their tags in order to filter out members who don't have a valid tag set
		ft := ""
		for _, Tag := range Tags {
			ft = tfield.Tag.Get(Tag)
			if ft != "" {
				break
			}
		}
		if ft == "" {
			continue
		}
		vfField := vf.Elem().Field(i)
		switch vfField.Type().String() {
		case "string":
			vfField.SetString(v.GetString(ft))
		case "int", "int64":
			vfField.SetInt(v.GetInt64(ft))
		case "float64":
			vfField.SetFloat(v.GetFloat64(ft))
		default:
		}
	}
	return object, nil
}
