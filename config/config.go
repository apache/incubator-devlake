package config

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Lowcase V for private this. You can use it by call GetConfig.
var v *viper.Viper = nil

func GetConfig() *viper.Viper {
	return v
}

// Set a struct for viper
// `Tags` represent the fields when setting config, and the fields in Tags shall prevail. `Tags` that appear first have higher priority.
func SetStruct(S interface{}, Tags ...string) error {
	err := validator.New().Struct(S)
	if err != nil {
		return err
	}

	v := GetConfig()
	tf := reflect.TypeOf(S)
	vf := reflect.ValueOf(S)

	for i := 0; i < tf.NumField(); i++ {
		tfield := tf.Field(i)
		vfield := vf.Field(i)

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

		v.Set(ft, vfield.Interface())
	}
	return v.WriteConfig()
}

// Set default value for no .env or .env not set it
func setDefaultValue() {
	v.SetDefault("PORT", ":8080")
	v.SetDefault("PLUGIN_DIR", "bin/plugins")
	v.SetDefault("TEMPORAL_TASK_QUEUE", "DEVLAKE_TASK_QUEUE")
}

func init() {
	// create the object and load the .env file
	v = viper.New()
	v.SetConfigFile(".env")
	err := v.ReadInConfig()
	if err != nil {
		logrus.Warn("Failed to read [.env] file:", err)
	}
	v.AutomaticEnv()

	setDefaultValue()
	// This line is essential for reading and writing
	v.WatchConfig()
}
