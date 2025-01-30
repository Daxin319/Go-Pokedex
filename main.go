package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	supportedCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
	supportedCommands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp(supportedCommands),
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if command, ok := supportedCommands[input[0]]; ok {
			if err := command.callback(); err != nil {
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
	callback    func() error
}

func cleanInput(text string) []string {
	split := strings.Fields(strings.ToLower(text))

	return split

}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(input map[string]cliCommand) func() error {
	return func() error {
		fmt.Print("Welcome to the Pokedex!\nUsage:\n\n\n")

		for _, command := range input {
			fmt.Printf("%q: %q\n", input[command.name].name, input[command.name].description)
		}

		return nil
	}
}
