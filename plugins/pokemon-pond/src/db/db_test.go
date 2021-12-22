package db

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetDBClient(t *testing.T) {
	uri := "postgres://postgres:postgresWhat@localhost:5432/lake?sslmode=disable"

	err := GetDBClient(uri)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.NotEmpty(t, DBClient)

}

func TestPushItem(t *testing.T) {
	uri := "postgres://postgres:postgresWhat@localhost:5432/lake?sslmode=disable"

	err := GetDBClient(uri)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer DBClient.Close()

	pokemon := PokemonDB{
		ID:     1,
		Cost:   1000,
		Detail: "test",
	}
	// delete test data
	r := DBClient.Delete(&pokemon)
	assert.NoError(t, r.Error)

	// insert
	r = DBClient.Create(&pokemon)
	assert.NoError(t, r.Error)

	// select
	var results []PokemonDB
	DBClient.Table(DBTABLE).Find(&results, "id = ?", 1)
	assert.Equal(t, 1000, results[0].Cost)

}

func TestGetTotalCost(t *testing.T) {
	err := GetDBClient(DBURI)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer DBClient.Close()

	// delete test data
	r := DBClient.Delete(&PokemonDB{})
	assert.NoError(t, r.Error)

	pokemons := []PokemonDB{
		{ID: 1, Cost: 1000, Detail: "test1"},
		{ID: 2, Cost: 2000, Detail: "test2"},
	}

	// insert
	for _, p := range pokemons {
		r := DBClient.Create(&p)
		assert.NoError(t, r.Error)
	}

	var itemIds = []int{1, 2}
	result, err := GetTotalCost(itemIds)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	assert.Equal(t, 3000, result)

}
