package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	apilogic "github.com/Daxin319/Go-Pokedex/apiLogic"
	pokecache "github.com/Daxin319/Go-Pokedex/internal"
)

func main() {
	// create cache with 5 second memory
	cache := pokecache.NewCache(5 * time.Second)

	// initialize config with zero values
	config := &apilogic.Config{
		NextURL:     "",
		PreviousURL: "",
		Area:        "",
		Pokemon:     "",
	}

	// map of caught pokemon
	caught := make(map[string]apilogic.Pokemon)

	// a map of supported commands for the pokedex
	supportedCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next of 20 regions in the Pokemon world",
			callback:    apilogic.CommandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 regions in the Pokemon world",
			callback:    apilogic.CommandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Allows you to see a list of pokemon found in your selected area with the syntax `explore name_of_region`",
			callback:    apilogic.CommandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon with the syntax `catch pokemon-name`",
			callback:    apilogic.CommandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Look at the information of a pokemon that you have captured.",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "check your total number of caught pokemon, and see a list of all caught pokemon",
			callback:    commandPokedex,
		},
	}
	// Had to add this one seperately because it's self-referential
	supportedCommands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp(supportedCommands, config),
	}

	//initialize scanner to wait for user input
	scanner := bufio.NewScanner(os.Stdin)

	//main loop
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		//"sanitize" the input by making it lowercase and separating by white space
		input := cleanInput(scanner.Text())
		//if the command is valid, do commands
		if command, ok := supportedCommands[input[0]]; ok {
			if len(input) > 1 {
				switch input[0] {
				case "explore":
					config.Area = input[1]
				case "catch", "inspect":
					config.Pokemon = input[1]
				default:
					continue
				}
			}
			if err := command.callback(cache, &caught, config); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown Command")
		}
	}
}

// Under conSTRUCTion
type cliCommand struct {
	name        string
	description string
	callback    func(cache *pokecache.Cache, m *map[string]apilogic.Pokemon, c *apilogic.Config) error
}

// Functionland

// do I have to explain this one?
func cleanInput(text string) []string {
	split := strings.Fields(strings.ToLower(text))

	return split

}

// exit the program
func commandExit(_ *pokecache.Cache, _ *map[string]apilogic.Pokemon, _ *apilogic.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// help command, lists all valid commands and brief description of each
func commandHelp(input map[string]cliCommand, _ *apilogic.Config) func(_ *pokecache.Cache, _ *map[string]apilogic.Pokemon, c *apilogic.Config) error {
	return func(_ *pokecache.Cache, _ *map[string]apilogic.Pokemon, c *apilogic.Config) error {
		fmt.Print("Welcome to the Pokedex!\nUsage:\n\n\n")

		for _, command := range input {
			fmt.Printf("%q: %q\n", input[command.name].name, input[command.name].description)
		}

		return nil
	}
}

// inspect command, lists information on caught pokemon
func commandInspect(_ *pokecache.Cache, m *map[string]apilogic.Pokemon, c *apilogic.Config) error {
	// set the key for the information lookup
	key := "https://pokeapi.co/api/v2/pokemon/" + c.Pokemon + "/"
	//check map for key and return data, return error if data does not exist
	if data, ok := (*m)[key]; ok {
		fmt.Printf("Name: %s\n", data.Name)
		fmt.Printf("Height: %d\n", data.Height)
		fmt.Printf("Weight: %d\n", data.Weight)
		fmt.Printf("Stats:\n")
		for _, stat := range data.Stats {
			fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _, element := range data.Types {
			fmt.Printf("  -%s\n", element.Type.Name)
		}
		return nil
	} else {
		return fmt.Errorf("you haven't caught that pokemon")
	}
}

// pokedex command, prints a list of all caught pokemon and a total count
func commandPokedex(_ *pokecache.Cache, m *map[string]apilogic.Pokemon, _ *apilogic.Config) error {
	fmt.Printf("So far, you have caught %d pokemon!\n\n", len(*m))
	fmt.Printf("Captured pokemon:\n")
	for _, pokemon := range *m { // this doesn't need parens but (*m)[key] does? ugh
		fmt.Printf("  -%s\n", pokemon.Name)
	}
	return nil
}
