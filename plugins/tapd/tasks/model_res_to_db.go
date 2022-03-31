package tasks

import (
	"fmt"
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
	//uint64Ins := uint64(1)
	for i := 0; i < tBd.NumField(); i++ {
		field := tBd.Field(i)
		//s := reflect.TypeOf(reflect.ValueOf(db).Elem().FieldByName(field.Name))

		if field.Type == reflect.TypeOf(uint64(1)) {
			if !reflect.ValueOf(body).Elem().FieldByName(field.Name).IsValid() {
				continue
			} else {
				v := reflect.ValueOf(body).Elem().FieldByName(field.Name)
				s := v.String()
				finalV, err := AtoIIgnoreEmpty(s)
				if err != nil {
					return nil, err
				}
				reflect.ValueOf(db).Elem().FieldByName(field.Name).
					Set(reflect.ValueOf(uint64(finalV)))
				fmt.Println(s)
			}

		} else if field.Type == reflect.TypeOf(&timeTest) {
			v := reflect.ValueOf(body).Elem().FieldByName(field.Name)
			s := v.String()
			finalV, err := core.ConvertStringToTimePtr(s)
			if err != nil {
				return nil, err
			}
			reflect.ValueOf(db).Elem().FieldByName(field.Name).
				Set(reflect.ValueOf(finalV))
		} else if field.Type == reflect.TypeOf(1) {
			v := reflect.ValueOf(body).Elem().FieldByName(field.Name)
			s := v.String()
			finalV, err := AtoIIgnoreEmpty(s)
			if err != nil {
				return nil, err
			}
			reflect.ValueOf(db).Elem().FieldByName(field.Name).
				Set(reflect.ValueOf(finalV))
		} else {
			if !reflect.ValueOf(body).Elem().FieldByName(field.Name).IsValid() {
				continue
			} else {
				reflect.ValueOf(db).Elem().FieldByName(field.Name).
					Set(reflect.ValueOf(body).Elem().FieldByName(field.Name))
			}
		}
		//if reflect.TypeOf(reflect.ValueOf(db).Elem().FieldByName(field.Name)). == "uint64" {
		//	fmt.Println("OK")
		//}
		//if reflect.TypeOf(reflect.ValueOf(db).Elem().FieldByName(field.Name)) ==  {
		//	s, err := AtoIIgnoreEmpty(reflect.ValueOf(db).Elem().FieldByName(field.Name).String())
		//	if err != nil {
		//		return nil, err
		//	}
		//	reflect.ValueOf(db).Elem().FieldByName(field.Name).Set(reflect.ValueOf(s))
		//}
		fmt.Println(field)
		//if field.Type == reflect.TypeOf(*time.Time) {
		//
		//}
		//reflect.ValueOf(db).Elem().FieldByName(field.Name).
		//	Set()

	}
	return db, nil
}
