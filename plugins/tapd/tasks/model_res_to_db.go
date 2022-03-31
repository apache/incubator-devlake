package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"reflect"
	"time"
)

func ResToDb(body interface{}, db interface{}) (interface{}, error) {
	vBd := reflect.ValueOf(db)
	if vBd.Kind() != reflect.Ptr {
		panic("entityPtr is not a pointer")
	}
	tBd := reflect.Indirect(vBd).Type()
	timeTest := time.Now()
	for i := 0; i < tBd.NumField(); i++ {
		field := tBd.Field(i)
		if !reflect.ValueOf(body).Elem().FieldByName(field.Name).IsValid() {
			continue
		}
		v := reflect.ValueOf(body).Elem().FieldByName(field.Name)
		s := v.String()
		switch field.Type {
		case reflect.TypeOf(uint64(1)):
			finalV, err := AtoIIgnoreEmpty(s)
			if err != nil {
				return nil, err
			}
			reflect.ValueOf(db).Elem().FieldByName(field.Name).
				Set(reflect.ValueOf(uint64(finalV)))
			break
		case reflect.TypeOf(&timeTest):
			finalV, err := core.ConvertStringToTimePtr(s)
			if err != nil {
				return nil, err
			}
			reflect.ValueOf(db).Elem().FieldByName(field.Name).
				Set(reflect.ValueOf(finalV))
			break
		case reflect.TypeOf(1):
			finalV, err := AtoIIgnoreEmpty(s)
			if err != nil {
				return nil, err
			}
			reflect.ValueOf(db).Elem().FieldByName(field.Name).
				Set(reflect.ValueOf(finalV))
			break
		default:
			reflect.ValueOf(db).Elem().FieldByName(field.Name).
				Set(v)
		}
	}
	return db, nil
}
