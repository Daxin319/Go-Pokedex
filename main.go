package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	apilogic "github.com/Daxin319/Go-Pokedex/apiLogic"
)

func main() {
	config := &apilogic.Config{
		NextURL:     "",
		PreviousURL: "",
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
			callback:    apilogic.commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 regions in the Pokemon world",
			callback:    apilogic.commandMapB,
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

type cliCommand struct {
	name        string
	description string
	callback    func(c *apilogic.Config) error
}

func cleanInput(text string) []string {
	split := strings.Fields(strings.ToLower(text))

	return split

}

func commandExit(c *apilogic.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(input map[string]cliCommand, c *apilogic.Config) func(c *apilogic.Config) error {
	return func(c *apilogic.Config) error {
		fmt.Print("Welcome to the Pokedex!\nUsage:\n\n\n")

		for _, command := range input {
			fmt.Printf("%q: %q\n", input[command.name].name, input[command.name].description)
		}

		return nil
	}
}
