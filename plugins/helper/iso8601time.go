package helper

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

/*
parse iso8601 datetime format for json.Unmarshal

declare your field type as:

type Foo struct {
	Created Iso8601Time
}

foo := &Foo{}
err := json.Unmarshal("{\"created\": \"2021-02-19T01:53:35.340+0800\"}", foo)
var time time.Time
time = foo.Created.ToTime()
*/

// TODO: move this to helper
type DateTimeFormatItem struct {
	Matcher *regexp.Regexp
	Format  string
}

var DateTimeFormats []DateTimeFormatItem

func init() {
	DateTimeFormats = []DateTimeFormatItem{
		{
			Matcher: regexp.MustCompile(`[+-][\d]{4}$`),
			Format:  "2006-01-02T15:04:05-0700",
		},
		{
			Matcher: regexp.MustCompile(`[+-][\d]{2}:[\d]{2}$`),
			Format:  "2006-01-02T15:04:05.000-07:00",
		},
	}
}

//type Iso8601Time time.Time
type Iso8601Time struct {
	time   time.Time
	format string
}

func (jt *Iso8601Time) String() string {
	format := jt.format
	if format == "" {
		format = DateTimeFormats[0].Format
	}
	return jt.time.Format(format)
}

func (jt Iso8601Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, jt.String())), nil
}

func (jt *Iso8601Time) UnmarshalJSON(b []byte) error {
	timeString := string(b)
	if timeString == "null" {
		return nil
	}
	timeString = strings.Trim(timeString, `"`)
	t, err := ConvertStringToTime(timeString)
	if err != nil {
		return err
	}
	jt.time = t
	return nil
}

func (jt *Iso8601Time) ToTime() time.Time {
	return jt.time
}

func (jt *Iso8601Time) ToNullableTime() *time.Time {
	if jt == nil {
		return nil
	}
	return &jt.time
}

func ConvertStringToTime(timeString string) (t time.Time, err error) {
	for _, formatItem := range DateTimeFormats {
		if formatItem.Matcher.MatchString(timeString) {
			t, err = time.Parse(formatItem.Format, timeString)
			return
		}
	}
	t, err = time.Parse(time.RFC3339, timeString)
	return
}

func Iso8601TimeToTime(iso8601Time *Iso8601Time) *time.Time {
	if iso8601Time == nil {
		return nil
	}
	t := iso8601Time.ToTime()
	return &t
}

// mapstructure.Decode with time.Time and Iso8601Time support
func DecodeMapStruct(input map[string]interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
				if t != reflect.TypeOf(Iso8601Time{}) && t != reflect.TypeOf(time.Time{}) {
					return data, nil
				}

				var tt time.Time
				var err error

				switch f.Kind() {
				case reflect.String:
					tt, err = ConvertStringToTime(data.(string))
				case reflect.Float64:
					tt = time.Unix(0, int64(data.(float64))*int64(time.Millisecond))
				case reflect.Int64:
					tt = time.Unix(0, data.(int64)*int64(time.Millisecond))
				}
				if err != nil {
					return data, nil
				}

				if t == reflect.TypeOf(Iso8601Time{}) {
					return Iso8601Time{time: tt}, nil
				}
				return tt, nil
			},
		),
		Result: result,
	})
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}
	return err
}
