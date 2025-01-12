package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/manuelgcsousa/pokevote/internal/data"
)

type pokeApiSpeciesData struct {
	Order int `json:"order"`
}

type pokeApiGenericDataResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type pokeApiGenericData struct {
	Results []pokeApiGenericDataResult `json:"results"`
}

type pokeApiSpecificData struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Order   int    `json:"order"`
	Species struct {
		Url string `json:"url"`
	} `json:"species"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

func main() {
	// Connect and init schema
	db, err := data.NewConnection()
	if err != nil {
		log.Fatal(err)
	}
	db.InitDatabase()

	// Fetch all pokemon
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/?limit=10000")
	if err != nil {
		log.Fatal(err)
	}

	var genericData pokeApiGenericData
	err = json.NewDecoder(resp.Body).Decode(&genericData)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch individual pokemon data and insert it into local DB
	for _, result := range genericData.Results {
		fetchSpecificPokeData(db, result)
	}
}

func fetchSpecificPokeData(db *data.Database, result pokeApiGenericDataResult) {
	resp, err := http.Get(result.URL)
	if err != nil {
		fmt.Println(err)
		return
	}

	var specificData pokeApiSpecificData
	err = json.NewDecoder(resp.Body).Decode(&specificData)
	if err != nil {
		fmt.Println(err)
	}

	newPokemon := &data.Pokemon{
		Id:        specificData.Id,
		DexId:     specificData.Order,
		Name:      strings.ReplaceAll(strings.Title(specificData.Name), "-", " "),
		SpriteUrl: specificData.Sprites.FrontDefault,
	}

	if specificData.Order == -1 {
		resp, err = http.Get(specificData.Species.Url)
		if err != nil {
			fmt.Println(err)
			return
		}

		var speciesData pokeApiSpeciesData
		err = json.NewDecoder(resp.Body).Decode(&speciesData)
		if err != nil {
			log.Fatal(err)
		}

		newPokemon.DexId = speciesData.Order
	}

	err = db.InsertPokemon(newPokemon)
	if err != nil {
		fmt.Println(err)
	}
}
