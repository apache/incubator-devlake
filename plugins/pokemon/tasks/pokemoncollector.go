package tasks

import (
	"net/url"
	"time"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"

	"github.com/merico-dev/lake/plugins/core"
	pok "github.com/merico-dev/lake/plugins/pokemon/models"
)

// getQuery - returns the query parameters given a full URL
func getQuery(fullURL string) (url.Values, error) {
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return map[string][]string{}, err
	}

	return url.ParseQuery(parsedURL.RawQuery)
}

func CollectAllPokemon() error {
	pokemonClient := core.NewApiClient(
		"https://pokeapi.co/api/v2/",
		nil,
		10*time.Second,
		3,
	)

	var query url.Values

	// loop through and get all the pokemon
	for {
		pokemonColl := &pok.PokemonCollection{}
		logger.Debug("Fetching pokemon with query params: ", query)

		res, err := pokemonClient.Get("/pokemon", &query, nil)
		if err != nil {
			return err
		}

		err = core.UnmarshalResponse(res, pokemonColl)
		if err != nil {
			logger.Error("Error unmarshalling pokemon response: ", err)
			return nil
		}

		for _, pkmon := range pokemonColl.Results {
			err = lakeModels.Db.Save(pkmon).Error
			if err != nil {
				logger.Error("Error saving pokemon: ", err)
			}
		}

		query, err = getQuery(pokemonColl.Next)

		// in the event of an error or when we have no 
		// query params, break out of the loop
		if err != nil || len(query) == 0 {
			break
		}

		// when we have no more next, break
		if pokemonColl.Next == "" {
			break
		}
		
		// primitive rate limitting so we don't hammer the API
		<-time.After(1 * time.Second)

	}

	return nil
}
