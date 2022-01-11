package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/merico-dev/lake/logger"
)

type ApiResponse struct {
	Count   int
	Results []Pokemon
	Next    string
}

type Pokemon struct {
	Name string
	Url  string
}

func UnmarshalResponse(res *http.Response, v interface{}) error {
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		logger.Print(fmt.Sprintf("UnmarshalResponse failed: %v\n%v\n\n", res.Request.URL.String(), string(resBody)))
		return err
	}
	return json.Unmarshal(resBody, &v)
}

func FetchPage(url string, c chan (bool)) {
	results, err := http.Get(url)
	if err != nil {
		fmt.Println("err1", err)
	}
	apiResponse := &ApiResponse{}
	err = UnmarshalResponse(results, apiResponse)

	if err != nil {
		fmt.Println("err2", err)
	}

	fmt.Println("JON >>> len(pokemons)", len(apiResponse.Results))

	for i, _ := range apiResponse.Results {
		// fmt.Println("JON >>> pok", pok)
		fmt.Println("JON >>> i", i)
	}
	c <- true
}

func main() {
	c := make(chan (bool))
	go FetchPage("https://pokeapi.co/api/v2/pokemon?offset=20&limit=20", c)
	go FetchPage("https://pokeapi.co/api/v2/pokemon?offset=40&limit=20", c)
	go FetchPage("https://pokeapi.co/api/v2/pokemon?offset=60&limit=20", c)
	go FetchPage("https://pokeapi.co/api/v2/pokemon?offset=80&limit=20", c)
	go FetchPage("https://pokeapi.co/api/v2/pokemon?offset=100&limit=20", c)
}
