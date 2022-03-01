package config

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Lowcase V for private this. You can use it by call GetConfig.
var v *viper.Viper = nil

func GetConfig() *viper.Viper {
	return v
}

// Set a struct for viper
// if IgnoreEmpty is set to true, empty string will be ignored
func SetStruct(S interface{}, IgnoreEmpty bool) {
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
		ft := tfield.Tag.Get("json")
		if ft == "" {
			ft = tfield.Tag.Get("mapstructure")
		}
		if ft == "" {
			continue
		}

		if IgnoreEmpty {
			switch tfield.Type.Kind() {
			case reflect.String:
				if vfield.String() == "" {
					continue
				}
			}
		}

		v.Set(ft, vfield.Interface())
	}
	v.WriteConfig()
}

// Set default value for no .env or .env not set it
func setDefaultValue() {
	v.SetDefault("PORT", ":8080")
	v.SetDefault("PLUGIN_DIR", "bin/plugins")
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
