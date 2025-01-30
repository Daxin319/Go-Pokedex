package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	config := &Config{
		nextURL:     "",
		previousURL: "",
	}
	supportedCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next of 20 regions in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 regions in the Pokemon world",
			callback:    commandMapB,
		},
	}
	supportedCommands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp(supportedCommands, config),
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if command, ok := supportedCommands[input[0]]; ok {
			if err := command.callback(config); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown Command")
		}
	}
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
type Config struct {
	nextURL     string
	previousURL string
}
type cliCommand struct {
	name        string
	description string
	callback    func(c *Config) error
}

func cleanInput(text string) []string {
	split := strings.Fields(strings.ToLower(text))

	return split

}

func commandExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(input map[string]cliCommand, c *Config) func(c *Config) error {
	return func(c *Config) error {
		fmt.Print("Welcome to the Pokedex!\nUsage:\n\n\n")

		for _, command := range input {
			fmt.Printf("%q: %q\n", input[command.name].name, input[command.name].description)
		}

		return nil
	}
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
