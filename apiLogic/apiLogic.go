package apilogic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
type Config struct {
	NextURL     string
	PreviousURL string
}

func CommandMap(c *Config) error {
	var url string
	if c.NextURL != "" {
		url = c.NextURL
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
		c.PreviousURL = *locations.Previous
	} else {
		c.PreviousURL = url
	}
	c.NextURL = locations.Next
	return nil
}
func CommandMapB(c *Config) error {
	var url string
	if c.PreviousURL == "" || c.PreviousURL == "https://pokeapi.co/api/v2/location-area/" {
		fmt.Println("You are on the first page")
		c.NextURL = "https://pokeapi.co/api/v2/location-area/"
		return nil
	} else {
		url = c.PreviousURL
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
	c.NextURL = locations.Next
	if locations.Previous != nil {
		c.PreviousURL = *locations.Previous
	} else {
		c.PreviousURL = ""
	}
	return nil
}
