package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo string

func (f *Foo) Description() string {
	return "foo"
}

func (f *Foo) Init() {
}

func (f *Foo) Execute(options map[string]interface{}, taskId uint64, progress chan<- float32) {
	close(progress)
}

func (f *Foo) RootPkgPath() string {
	return "path/to/foo"
}

type Bar string

func (b *Bar) Description() string {
	return "foo"
}

func (b *Bar) Init() {
}

func (b *Bar) Execute(options map[string]interface{}, taskId uint64, progress chan<- float32) {
	close(progress)
}

func (b *Bar) RootPkgPath() string {
	return "path/to/bar"
}

func TestHub(t *testing.T) {
	var foo Foo
	assert.Nil(t, RegisterPlugin("foo", &foo))
	var bar Bar
	assert.Nil(t, RegisterPlugin("bar", &bar))

	f, _ := GetPlugin("foo")
	assert.Equal(t, &foo, f)

	fn, _ := FindPluginNameBySubPkgPath("path/to/foo/models")
	assert.Equal(t, fn, "foo")

	b, _ := GetPlugin("bar")
	assert.Equal(t, &bar, b)

	bn, _ := FindPluginNameBySubPkgPath("path/to/bar/models")
	assert.Equal(t, bn, "bar")
}
