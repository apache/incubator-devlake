package okgen

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/merico-dev/lake/plugins/core"
)

type OriginKeyGenerator struct {
	prefix  string
	pkNames []string
	pkTypes []reflect.Type
}

func walkFields(t reflect.Type, pkNames *[]string, pkTypes *[]reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			walkFields(field.Type, pkNames, pkTypes)
		} else {
			gormTag := field.Tag.Get("gorm")

			// TODO: regex?
			if gormTag != "" && strings.Contains(strings.ToLower(gormTag), "primarykey") {
				*pkNames = append(*pkNames, field.Name)
				*pkTypes = append(*pkTypes, field.Type)
			}
		}
	}
}

func NewOriginKeyGenerator(entityPtr interface{}) *OriginKeyGenerator {
	v := reflect.ValueOf(entityPtr)
	if v.Kind() != reflect.Ptr {
		panic("entityPtr is not a pointer")
	}
	t := reflect.Indirect(v).Type()

	// find out which plugin holds the entity
	pluginName, err := core.FindPluginNameBySubPkgPath(t.PkgPath())
	if err != nil {
		panic(err)
	}
	// find out entity type name
	structName := t.Name()

	// find out all primkary keys and their types
	pkNames := make([]string, 0, 1)
	pkTypes := make([]reflect.Type, 0, 1)

	walkFields(t, &pkNames, &pkTypes)

	if len(pkNames) == 0 {
		panic(fmt.Errorf("no primary key found for %s:%s", pluginName, structName))
	}

	return &OriginKeyGenerator{
		prefix:  fmt.Sprintf("%s:%s", pluginName, structName),
		pkNames: pkNames,
		pkTypes: pkTypes,
	}
}

func (g *OriginKeyGenerator) Generate(pkValues ...interface{}) string {
	originKey := g.prefix
	for i, pkValue := range pkValues {
		// type checking
		if reflect.ValueOf(pkValue).Type() != g.pkTypes[i] {
			panic(fmt.Errorf("primary key type does not match: %s should be %s",
				g.pkNames[i],
				g.pkTypes[i].Name(),
			))
		}
		// append pk
		originKey += ":" + fmt.Sprintf("%v", pkValue)
	}
	return originKey
}
