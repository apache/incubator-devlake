package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/lake/src/plugins/pokemon-pond/src/db"
)

// Meta ...
type Meta struct {
	Count int `json:"count"`
}

// Pokemon ...
type Pokemon struct {
	URL               string            `json:"url,omitempty"`
	ID                int               `json:"id"` // PK
	Attributes        []attribute       `json:"attributes"`
	BabyTriggerFor    string            `json:"baby_trigger_for"` // TODO(Ted): need to check data. Not sure what type it is.
	Category          category          `json:"category"`
	Cost              int               `json:"cost"`
	EffectEntries     []effectentry     `json:"effect_entries"`
	FlavorTextEntries []flavorTextEntry `json:"flavor_text_entries"`
	//FlingEffect       string            `json:"fling_effect"` // TODO(Ted): need to check data. Not sure what type it is.
	FlingPower  int          `json:"fling_power"`
	GameIndices []gameIndice `json:"game_indices"`
	Name        string       `json:"name"`
	Names       []name       `json:"names"`
	Sprites     sprite       `json:"sprites"`
	// Machines // TODO(Ted): need to check
}

type attribute struct {
	nameAndURL
}

type category struct {
	nameAndURL
}

type effectentry struct {
	Effect      string   `json:"effect"`
	Language    language `json:"language"`
	ShortEffect string   `json:"short_effect"`
}

type flavorTextEntry struct {
	language     `json:"language"`
	Text         string     `json:"text"`
	VersionGroup nameAndURL `json:"version_group"`
}

type gameIndice struct {
	GameIndex  int        `json:"game_index"`
	Generation nameAndURL `json:"generation"`
}

type language struct {
	nameAndURL
}

type name struct {
	Name     string `json:"name"`
	language `json:"language"`
}

type sprite struct {
	Default string `json:"default"`
}

type nameAndURL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PokeURL ...
var PokeURL = "https://pokeapi.co/api/v2/pokemon"  // TODO(Ted): need to set to read from config or something in real world
var fetchURL = "https://pokeapi.co/api/v2/item/%s" // TODO(Ted): need to set to read from config or something in real world

// NewPokemon is to make Pokemon object
func NewPokemon() *Pokemon {
	return &Pokemon{}
}

// httpGet module
func httpGet(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetTotalCount is to get total fetch count
func GetTotalCount(url string) (int, error) {
	resp, err := httpGet(url)
	if err != nil {
		return 0, err
	}

	if resp == nil {
		return 0, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Fail GET %s. StatusCode: %d", url, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var m Meta
	err = json.Unmarshal(body, &m)
	if err != nil {
		return 0, err
	}

	return m.Count, nil
}

// Run is to start to fether and push item
func Run(totalCount int, itemID chan int) {
	for i := 1; i <= totalCount; i++ {
		itemID <- i
	}
	close(itemID)
}

// Fetcher is to get Pokemon item information
func Fetcher(itemID chan int, fetch chan Pokemon, prrError *map[int]error) {
	for i := range itemID {
		url := fmt.Sprintf(fetchURL, strconv.Itoa(i))
		resp, err := httpGet(url)
		if err != nil {
			(*prrError)[i] = err
		}

		if resp == nil {
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			(*prrError)[i] = err
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			(*prrError)[i] = err
			continue
		}

		var p Pokemon
		err = json.Unmarshal(body, &p)
		if err != nil {
			(*prrError)[i] = err
			continue
		}

		if i != p.ID {
			log.Printf("Different ID. index: %d, result itemID %d", i, p.ID)
		}

		fetch <- p
	}

	close(fetch)
	log.Printf("Fetcher done.")

}

// PushItem is to store pokemon data into PostgreDB
func PushItem(fetch chan Pokemon, done chan bool, prrError *map[int]error) {
	for p := range fetch {

		var pdb db.PokemonDB
		pdb.ID = p.ID
		pdb.Cost = p.Cost
		bp, err := json.Marshal(p)
		if err != nil {
			(*prrError)[p.ID] = err
		}
		pdb.Detail = string(bp)

		r := db.DBClient.Create(&pdb)

		// r := db.DBClient.Exec(
		// 	"INSERT INTO pokemon_dbs(id, cost, detail, created_at) VALUES (?,?,?,?)", pdb.ID, pdb.Cost, pdb.Detail, time.Now(),
		// )

		if r.Error != nil {
			(*prrError)[p.ID] = err
		}

	}

	done <- true
	log.Printf("Push done.")

}
