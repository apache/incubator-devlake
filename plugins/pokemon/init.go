package main

import (
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	pokemonModels "github.com/merico-dev/lake/plugins/pokemon/models"
)

func (p Pokemon) Init() {
	logger.Info("INFO >>> init go plugin", true)
	err := lakeModels.Db.AutoMigrate(
		&pokemonModels.Pokemon{},
	)
	if err != nil {
		logger.Error("Error migrating gitlab: ", err)
		panic(err)
	}
}
