package models

type Pokemon struct {
	PID        uint `gorm:"primary_key"`
	Name string `json:"name"`
	URL string `json:"url"`
}

type PokemonCollection struct {
	Count int
	Next string
	Previous string
	Results []*Pokemon
}