package core

type Plugin interface {
	Description() string
	Init()
	Execute(options map[string]interface{}, progress chan<- float32)
}
