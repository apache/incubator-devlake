package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	pokemonModels "github.com/merico-dev/lake/plugins/pokemon/models"
	"gorm.io/gorm/clause"
)

const pokemonURL = "https://pokeapi.co/api/v2/pokemon"

type Pokemon string

func (Pokemon) Description() string {
	return "To collect and enrich data from pokemon api"
}

func (Pokemon) Execute(options map[string]interface{}, p chan<- float32, ctx context.Context) {
	logger.Print("start pokemon plugin execution")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		pokemonURL,
		nil)
	if err != nil {
		logger.Error("error creating request", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("could not request", fmt.Errorf("%s %w", pokemonURL, err))
		return
	}
	defer res.Body.Close()

	r := PokemonResponse{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logger.Error("could not decode body", err)
		return
	}

	logger.Info("pokemon count: ", len(r.Results))

	pp := []*pokemonModels.Pokemon{}
	for _, p := range r.Results {
		pp = append(pp, &pokemonModels.Pokemon{
			Name: p.Name,
			URL:  p.URL,
		})
	}

	pchan := make(chan *pokemonModels.Pokemon)
	go func() {
		defer close(pchan)
		for i, p := range pp {
			pchan <- p
			// wait a minute after we send 10 for an hypothetical
			// 10 per 1 minute rate limit
			// there's probably better ways to achieve this.
			if i%10 == 0 {
				time.Sleep(time.Second * 60)
			}
		}
	}()

	// Fetch each pokemon with workers
	const maxWorkers = 10
	wg := sync.WaitGroup{}
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				// receive pokemon
				case p, ok := <-pchan:
					if !ok {
						return
					}
					_ = p
				// TODO: {lpf} do a request and fetch the
				// pokemon details for each pokemon and update
				// it in p
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	wg.Wait()

	// Store, could be batch inserted
	for _, p := range r.Results {
		err := lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&p).Error
		if err != nil {
			logger.Error("could not upsert: ", err)
			return
		}
	}
}

func (p Pokemon) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/gitlab"
}

func (p Pokemon) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return make(map[string]map[string]core.ApiResourceHandler)
}

type PokemonResponse struct {
	Count   int            `json:"count"`
	Next    string         `json:"next"`
	Results []*PokemonItem `json:"results"`
}

type PokemonItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var PluginEntry Pokemon
