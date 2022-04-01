package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"reflect"
)

// VoToDTO will transfer interface with all string field
// to another interface has similar fields but not same types
func VoToDTO(vo interface{}, dto interface{}) (interface{}, error) {
	vDto := reflect.ValueOf(dto)
	if vDto.Kind() != reflect.Ptr {
		panic("entityPtr is not a pointer")
	}
	tDto := reflect.Indirect(vDto).Type()
	for i := 0; i < tDto.NumField(); i++ {
		field := tDto.Field(i)
		if !reflect.ValueOf(vo).Elem().FieldByName(field.Name).IsValid() {
			continue
		}
		v := reflect.ValueOf(vo).Elem().FieldByName(field.Name)
		err := populate(reflect.ValueOf(dto).Elem().FieldByName(field.Name), v)
		if err != nil {
			return nil, err
		}
	}
	return dto, nil
}

//!+populate
func populate(v reflect.Value, value reflect.Value) error {
	switch v.Type().String() {
	case "uint64":
		finalV, err := AtoIIgnoreEmpty(value.String())
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(uint64(finalV)))
		break
	case "*time.Time":
		finalV, err := core.ConvertStringToTimePtr(value.String())
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(finalV))
		break
	case "int":
		finalV, err := AtoIIgnoreEmpty(value.String())
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(finalV))
		break
	default:
		v.Set(value)
	}
	return nil
}
