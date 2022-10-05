package logger

import (
	"fmt"
)

type Symbol string

// Symbols struct contains all symbols
type Symbols struct {
	Debug   Symbol
	Info    Symbol
	Warning Symbol
	Warn    Symbol
	Error   Symbol
	Fatal   Symbol
	Success Symbol
}

var normal = Symbols{
	Debug:   Symbol("λ"),
	Info:    Symbol("ℹ"),
	Success: Symbol("✔"),
	Warning: Symbol("⚠"),
	Warn:    Symbol("⚠"),
	Error:   Symbol("!!"),
	Fatal:   Symbol("✖"),
}

// String returns a printable representation of Symbols struct
func (s Symbols) String() string {
	return fmt.Sprintf("Debug: %s Info: %s Success: %s Warning: %s Error: %s Fatal: %s", s.Debug, s.Info, s.Success, s.Warning, s.Error, s.Fatal)
}
