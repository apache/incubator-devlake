package core

import (
	"fmt"
	"regexp"
	"strings"
	"time"
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
