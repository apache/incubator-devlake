package main

import (
	"log"
	"time"

	"github.com/lake/src/plugins/pokemon-pond/src/db"
	"github.com/lake/src/plugins/pokemon-pond/src/fetcher"
)

func main() {
	// Get DB conn

	log.Println("Start to run.... it does take a little bit time. wait please til finish")
	start := time.Now()

	err := db.GetDBClient(db.DBURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.DBClient.Close()

	// delete test data
	r := db.DBClient.Delete(&db.PokemonDB{})
	if r.Error != nil {
		log.Fatal(r.Error)
	}

	// Get Total count for fetch
	totalFetchCount, err := fetcher.GetTotalCount(fetcher.PokeURL)
	if err != nil {
		log.Fatal(err)
	}

	err = run(totalFetchCount)
	if err != nil {
		log.Fatal(err)
	}

	// Get totalCost
	// using eample item ids [251, 300]
	var itemIds = []int{251, 300} // itemid 251's cost is 1000. itemid 300's cost is 2000
	totalCost, err := db.GetTotalCost(itemIds)
	if err != nil {
		log.Printf("ERROR: Get Total Cost of IDs %v, err: %+v", itemIds, err)
	}

	log.Printf("Finish. Success. Total cost of Items[%v] is %d.", itemIds, totalCost)
	log.Printf("Elaspes time: %v", time.Since(start).Seconds())
}

func run(totalFetchCount int) error {

	done := make(chan bool)
	defer close(done)

	itemID := make(chan int)
	fetch := make(chan fetcher.Pokemon)
	prrError := make(map[int]error, 0)

	go fetcher.Run(totalFetchCount, itemID)
	go fetcher.Fetcher(itemID, fetch, &prrError)
	go fetcher.PushItem(fetch, done, &prrError)
	<-done

	if len(prrError) > 0 {
		// Do something with error. ex can retry or send alert ....
		log.Printf("%d items were not handled properly.\n", len(prrError))
		log.Println("Please see the prrError map for detailed error details.")
	}

	return nil

}
