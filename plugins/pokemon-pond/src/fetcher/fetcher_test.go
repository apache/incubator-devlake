package fetcher

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/lake/src/plugins/pokemon-pond/src/db"
	"github.com/stretchr/testify/assert"
)

func readSampleData(file string) (*string, error) {

	j, _ := os.Open(file)
	defer j.Close()

	byteValue, err := ioutil.ReadAll(j)
	if err != nil {
		return nil, err
	}
	sValue := string(byteValue)
	return &sValue, nil
}

func TestUnmarshal(t *testing.T) {

	f := "../../test/sample_251.json"
	data, err := readSampleData(f)
	if err != nil {
		t.Errorf("fail to read smaple data. %+v", err)
		t.FailNow()
	}

	pokenmon := NewPokemon()
	err = json.Unmarshal([]byte(*data), &pokenmon)
	if err != nil {
		t.Errorf("fail unmarshal. %+v", err)
		t.FailNow()
	}

	// I only compared some filed and count whether equal or not.
	// Because JSON is unodered. So compare this with hash or checksum or something...
	// It needs sort or something...

	// id
	expectedID := 251
	assert.Equal(t, expectedID, pokenmon.ID)

	// attribute
	expectedAttributeName := "holdable"
	assert.Equal(t, expectedAttributeName, pokenmon.Attributes[0].Name)

	// category
	expectedCategory := "species-specific"
	assert.Equal(t, expectedCategory, pokenmon.Category.Name)

	// cost
	expectedCost := 1000
	assert.Equal(t, expectedCost, pokenmon.Cost)

	// effect entries array
	expectedEffectentries0Effect := "Held by Ditto: Doubles the holder's initial Speed."
	assert.Equal(t, expectedEffectentries0Effect, pokenmon.EffectEntries[0].Effect)

	// flavor_text_entries
	expectedFlavorTextEntriesCount := 34
	assert.Equal(t, expectedFlavorTextEntriesCount, len(pokenmon.FlavorTextEntries))

	// TODO(Ted): add more

}

func TestGetTotalCount(t *testing.T) {
	url := "https://pokeapi.co/api/v2/pokemon"
	count, err := GetTotalCount(url)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	expectedCount := 1118
	assert.Equal(t, expectedCount, count)
}

func TestFetcher(t *testing.T) {
	f := "https://pokeapi.co/api/v2/item/50000"
	resp, err := httpGet(f)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, 404, resp.StatusCode)
}

func TestFetcherAndPushItem(t *testing.T) {

	err := db.GetDBClient(db.DBURI)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer db.DBClient.Close()

	r := db.DBClient.Delete(&db.PokemonDB{})
	assert.NoError(t, r.Error)

	//link := make(chan string)

	done := make(chan bool)
	defer close(done)

	itemID := make(chan int)
	fetch := make(chan Pokemon)
	prrError := make(map[int]error, 0)

	totalCount := 200

	go Run(totalCount, itemID)
	go Fetcher(itemID, fetch, &prrError)
	go PushItem(fetch, done, &prrError)
	<-done

	log.Printf("DEBUG ERROR: %+v", prrError)

}
