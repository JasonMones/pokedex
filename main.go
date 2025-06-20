package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/JasonMones/pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, args []string) error
}

type config struct {
	prev    string
	next    string
	cache   *pokecache.Cache
	pokedex map[string]pokemon
}

func generateRegistery() map[string]cliCommand {
	registry := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays 20 in-game map areas at a time",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 in-game map areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "lists the pokemon that can be found in a given area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "attempt to catch the specified pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "look at the stats of any caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "look at all pokemon caught",
			callback:    commandPokedex,
		},
	}
	return registry
}

func main() {
	regi := generateRegistery()
	c := config{
		prev:    "",
		next:    "https://pokeapi.co/api/v2/location-area/",
		cache:   pokecache.NewCache(20 * time.Second),
		pokedex: make(map[string]pokemon),
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		input := cleanInput(scanner.Text())
		inputtedCommand := input[0]

		command, exists := regi[inputtedCommand]
		if !exists {
			fmt.Println("unknown command")
			continue
		}

		if len(input) == 1 { //remove the command from the input
			input = make([]string, 0)
		} else {
			input = slices.Delete(input, 0, 1)
		}

		err := command.callback(&c, input)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	words := strings.Fields(lowerText)
	return words // first puts everything into lowercase, then
	// seperates every word by whitespaces into a []string
}
