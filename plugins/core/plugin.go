package core

type Plugin interface {
	Description() string
	Execute(options map[string]interface{}, progress chan<- float32)
}
