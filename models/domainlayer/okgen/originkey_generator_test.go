package okgen

import (
	"context"
	"testing"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type FooPlugin string

func (f *FooPlugin) Description() string {
	return "foo"
}

func (f *FooPlugin) Init() {
}

func (f *FooPlugin) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	return nil
}

func (f *FooPlugin) RootPkgPath() string {
	return "github.com/merico-dev/lake"
}

func (f *FooPlugin) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return make(map[string]map[string]core.ApiResourceHandler)
}

type FooModel struct {
	gorm.Model
}

func TestOriginKeyGenerator(t *testing.T) {
	var foo FooPlugin
	assert.Nil(t, core.RegisterPlugin("fooplugin", &foo))

	g := NewOriginKeyGenerator(&FooModel{})
	assert.Equal(t, g.prefix, "fooplugin:FooModel")

	originKey := g.Generate(uint(2))
	assert.Equal(t, "fooplugin:FooModel:2", originKey)

	assert.Panics(t, func() {
		NewOriginKeyGenerator(&foo)
	})

	assert.Panics(t, func() {
		g.Generate("asdf")
	})
}
