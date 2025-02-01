package apilogic

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	pokecache "github.com/Daxin319/Go-Pokedex/internal"
)

// conSTRUCTion zone
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
	Pokemon     string
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Forms          []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []any  `json:"past_abilities"`
	PastTypes     []any  `json:"past_types"`
	Species       struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

// conjunction junction, what's your function?

// map command, displays page 1 of the area list and paginates +1 each subsequent call
func CommandMap(cache *pokecache.Cache, _ *map[string]Pokemon, c *Config) error {
	var url string
	//set url for api pull
	if c.NextURL != "" {
		url = c.NextURL
	} else {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	//perform api pull, checkes cache for stored version before making a pull
	body, err := fetchDataWithCache(cache, url)
	if err != nil {
		return fmt.Errorf("error fetching data from api")
	}
	//initialize struct to unmarshal json
	locations := Locations{}
	//unmarshal the json
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return fmt.Errorf("error unmarshaling json file")
	}
	//iterate through the structs and print the name of each available location
	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	//set new destination for map and mapb
	if locations.Previous != nil {
		c.PreviousURL = *locations.Previous
	} else {
		c.PreviousURL = url
	}
	c.NextURL = locations.Next
	return nil
}

// mapb command, moves back 1 page in the area list, does not loop first to last page but will loop last to first
func CommandMapB(cache *pokecache.Cache, _ *map[string]Pokemon, c *Config) error {
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

// explore command, returns the list of available pokemon in a user specified area
func CommandExplore(cache *pokecache.Cache, _ *map[string]Pokemon, c *Config) error {
	//set the url for the api pull
	areaURL := "https://pokeapi.co/api/v2/location-area/" + c.Area + "/"
	//check cache for url and perform api pull if not present
	body, err := fetchDataWithCache(cache, areaURL)
	if err != nil {
		return fmt.Errorf("error fetching data from api")
	}
	//initialize struct to unmarshal json
	area := Area{}
	//unmarshal the json
	err = json.Unmarshal(body, &area)
	if err != nil {
		return fmt.Errorf("error unmarshaling json file")
	}
	//print message to user and iterate through structs to list available pokemon
	fmt.Printf("\nExploring %s!\n\n", c.Area)
	fmt.Printf("The following pokemon can be found in %s:\n\n", c.Area)
	for _, encounter := range area.PokemonEncounters {
		fmt.Println("- " + encounter.Pokemon.Name)
	}
	return nil
}

// capture command, attempts to catch a specified pokemon
func CommandCatch(cache *pokecache.Cache, m *map[string]Pokemon, c *Config) error {
	//set the url for the api pull
	pokemonURL := "https://pokeapi.co/api/v2/pokemon/" + c.Pokemon + "/"
	//check cache for url and perform pull if absent
	body, err := fetchDataWithCache(cache, pokemonURL)
	if err != nil {
		return fmt.Errorf("error fetching data from api")
	}
	//initialize struct to unmarshal json
	target := Pokemon{}
	//unmarshal the json
	err = json.Unmarshal(body, &target)
	if err != nil {
		return fmt.Errorf("error unmarshaling json file")
	}
	//print message to user and try to catch pokemon
	fmt.Printf("Throwing a Pokeball at %s...\n\n", target.Name)
	catchAttempt := rand.Intn(700)
	if catchAttempt >= target.BaseExperience {
		fmt.Printf("%s has been caught!\n", target.Name)
		fmt.Printf("Adding %s to list of caught pokemon!\n\n", target.Name)
		(*m)[pokemonURL] = target
	} else {
		fmt.Printf("Oh no! %s escaped!\n", target.Name)
		fmt.Printf("Catch roll: %d, Base XP: %d\n\n", catchAttempt, target.BaseExperience)
	}
	return nil
}

// perform api pull
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

// check cache to see if data exists, if so return data, else perform api pull
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
