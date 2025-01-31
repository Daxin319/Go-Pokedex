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
	}

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
				config.Area = input[1]
			}
			if err := command.callback(cache, config); err != nil {
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
	callback    func(cache *pokecache.Cache, c *apilogic.Config) error
}

// Functionland

// do I have to explain this one?
func cleanInput(text string) []string {
	split := strings.Fields(strings.ToLower(text))

	return split

}

// exit the program
func commandExit(_ *pokecache.Cache, _ *apilogic.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// help command, lists all valid commands and brief description of each
func commandHelp(input map[string]cliCommand, _ *apilogic.Config) func(_ *pokecache.Cache, c *apilogic.Config) error {
	return func(_ *pokecache.Cache, c *apilogic.Config) error {
		fmt.Print("Welcome to the Pokedex!\nUsage:\n\n\n")

		for _, command := range input {
			fmt.Printf("%q: %q\n", input[command.name].name, input[command.name].description)
		}

		return nil
	}
}
