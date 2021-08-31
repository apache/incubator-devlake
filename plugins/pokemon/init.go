package main

import (
	"github.com/merico-dev/lake/models"
	pok "github.com/merico-dev/lake/plugins/pokemon/models"
)

func (plugin Pokemon) Init() {
	err := models.Db.AutoMigrate(&pok.Pokemon{})
	if err != nil {
		panic(err)
	}
}
