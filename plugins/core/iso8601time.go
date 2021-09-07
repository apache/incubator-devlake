package core

import (
	"fmt"
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
type Iso8601Time time.Time

const ISO_8601_FORMAT = "2006-01-02T15:04:05-0700"

func (jt *Iso8601Time) String() string {
	t := time.Time(*jt)
	return t.Format(ISO_8601_FORMAT)
}

func (jt Iso8601Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, jt.String())), nil
}

func (jt *Iso8601Time) UnmarshalJSON(b []byte) error {
	timeString := strings.Trim(string(b), `"`)
	if strings.ToLower(timeString) == "null" {
		return nil
	}
	t, err := time.Parse(ISO_8601_FORMAT, timeString)
	if err == nil {
		*jt = Iso8601Time(t)
		return nil
	}
	return fmt.Errorf("invalid date format: %s", timeString)
}

func (jt *Iso8601Time) ToTime() time.Time {
	return time.Time(*jt)
}
