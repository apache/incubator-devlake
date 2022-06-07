package core

type ConfigGetter interface {
	GetString(name string) string
}

type InjectConfigGetter interface {
	SetConfigGetter(getter ConfigGetter)
}
