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
	cache := pokecache.NewCache(5 * time.Second)
	config := &apilogic.Config{
		NextURL:     "",
		PreviousURL: "",
		Area:        "",
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

type cliCommand struct {
	name        string
	description string
	callback    func(cache *pokecache.Cache, c *apilogic.Config) error
}

func cleanInput(text string) []string {
	split := strings.Fields(strings.ToLower(text))

	return split

}

func commandExit(_ *pokecache.Cache, _ *apilogic.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(input map[string]cliCommand, _ *apilogic.Config) func(_ *pokecache.Cache, c *apilogic.Config) error {
	return func(_ *pokecache.Cache, c *apilogic.Config) error {
		fmt.Print("Welcome to the Pokedex!\nUsage:\n\n\n")

		for _, command := range input {
			fmt.Printf("%q: %q\n", input[command.name].name, input[command.name].description)
		}

		return nil
	}
}
