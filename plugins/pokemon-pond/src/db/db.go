package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

// PokemonDB ...
// gorm.Model definition
type PokemonDB struct {
	ID        int    `gorm:"primaryKey"`
	Cost      int    // cost
	Detail    string // store whole pokemon into by string of json
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DBClient ...
var DBClient *gorm.DB

// DBURI ...
var DBURI = "postgres://postgres:postgresWhat@localhost:5432/lake?sslmode=disable" // TODO(Ted): need to set to read from config or something in real world

// DBTABLE ...
var DBTABLE = "pokemon_dbs" // TODO(Ted): need to set to read from config or something in real world

// GetDBClient is to make db conn
// It has been using gorm
func GetDBClient(uri string) error {
	if uri == "" {
		return errors.New("you must set your 'POSTGRE_URI' environmental variable")
	}

	db, err := gorm.Open("postgres", uri)
	if err != nil {
		return fmt.Errorf("Unable to connect to database: %v", err)
	}

	db.DB().SetMaxOpenConns(25)
	db.DB().SetMaxIdleConns(25)
	db.DB().SetConnMaxLifetime(5 * time.Minute)

	DBClient = db
	DBClient.AutoMigrate(&PokemonDB{})

	return nil
}

// GetTotalCost is to get total cost of items what you want
func GetTotalCost(itemIds []int) (int, error) {
	if len(itemIds) == 0 {
		return 0, nil
	}

	err := GetDBClient(DBURI)
	if err != nil {
		return 0, err
	}
	//var totalCost []int

	type TotalCost struct {
		Cost int
	}
	var totalCost TotalCost

	DBClient.Raw("SELECT SUM(cost) as cost FROM pokemon_dbs WHERE id IN (?)", itemIds).Scan(&totalCost)

	return totalCost.Cost, nil

}
