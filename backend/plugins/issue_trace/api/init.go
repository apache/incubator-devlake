package api

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/go-playground/validator/v10"
)

var BasicRes context.BasicRes
var Validator *validator.Validate

func Init(basicRes context.BasicRes) {
	BasicRes = basicRes
	Validator = validator.New()
}
