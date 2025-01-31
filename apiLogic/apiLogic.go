package apilogic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	pokecache "github.com/Daxin319/Go-Pokedex/internal"
)

type Locations struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Area struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Config struct {
	NextURL     string
	PreviousURL string
	Area        string
}

func CommandMap(cache *pokecache.Cache, c *Config) error {
	var url string
	if c.NextURL != "" {
		url = c.NextURL
	} else {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	body, err := fetchDataWithCache(cache, url)
	if err != nil {
		return fmt.Errorf("error fetching data from api")
	}
	locations := Locations{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return fmt.Errorf("error unmarshaling json file")
	}
	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	if locations.Previous != nil {
		c.PreviousURL = *locations.Previous
	} else {
		c.PreviousURL = url
	}
	c.NextURL = locations.Next
	return nil
}
func CommandMapB(cache *pokecache.Cache, c *Config) error {
	var url string
	if c.PreviousURL == "" || c.PreviousURL == "https://pokeapi.co/api/v2/location-area/" {
		fmt.Println("You are on the first page")
		c.NextURL = "https://pokeapi.co/api/v2/location-area/"
		return nil
	} else {
		url = c.PreviousURL
	}
	body, err := fetchDataWithCache(cache, url)
	if err != nil {
		return fmt.Errorf("error fetching data from api")
	}
	locations := Locations{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return fmt.Errorf("error unmarshaling json file")
	}
	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	c.NextURL = locations.Next
	if locations.Previous != nil {
		c.PreviousURL = *locations.Previous
	} else {
		c.PreviousURL = ""
	}
	return nil
}

func CommandExplore(cache *pokecache.Cache, c *Config) error {
	areaURL := "https://pokeapi.co/api/v2/location-area/" + c.Area + "/"
	body, err := fetchDataWithCache(cache, areaURL)
	if err != nil {
		return fmt.Errorf("error fetching data from api")
	}
	area := Area{}
	err = json.Unmarshal(body, &area)
	if err != nil {
		return fmt.Errorf("error unmarshaling json file")
	}
	fmt.Printf("\nExploring %s!\n\n", c.Area)
	fmt.Printf("The following pokemon can be found in %s:\n\n", c.Area)
	for _, encounter := range area.PokemonEncounters {
		fmt.Println("- " + encounter.Pokemon.Name)
	}
	return nil
}

func fetchDataFromAPI(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error pulling location data from pokeapi")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body")
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with error code %d and \nbody: %s", resp.StatusCode, resp.Body)
	}
	return body, nil

}

func fetchDataWithCache(cache *pokecache.Cache, url string) ([]byte, error) {
	if data, ok := cache.Get(url); ok {
		return data, nil
	}
	data, err := fetchDataFromAPI(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data from API")
	}
	cache.Add(url, data)
	return data, nil
}
