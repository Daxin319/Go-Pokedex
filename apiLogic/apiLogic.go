package apilogic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	nextURL     string
	previousURL string
}

type Locations struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMap(c *Config) error {
	var url string
	if c.nextURL != "" {
		url = c.nextURL
	} else {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error pulling location data from pokeapi")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("error reading response body")
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("response failed with error code %d and \nbody: %s", resp.StatusCode, resp.Body)
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
		c.previousURL = *locations.Previous
	} else {
		c.previousURL = url
	}
	c.nextURL = locations.Next

	return nil
}

func commandMapB(c *Config) error {
	var url string
	if c.previousURL == "" || c.previousURL == "https://pokeapi.co/api/v2/location-area/" {
		fmt.Println("You are on the first page")
		c.nextURL = "https://pokeapi.co/api/v2/location-area/"
		return nil
	} else {
		url = c.previousURL
	}
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error pulling location data from pokeapi")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("error reading response body")
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("response failed with error code %d and \nbody: %s", resp.StatusCode, resp.Body)
	}

	locations := Locations{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return fmt.Errorf("error unmarshaling json file")
	}

	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}

	c.nextURL = locations.Next
	if locations.Previous != nil {
		c.previousURL = *locations.Previous
	} else {
		c.previousURL = ""
	}

	return nil
}
