package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/JasonMones/pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

type config struct {
	prev  string
	next  string
	cache *pokecache.Cache
}

type location_batch struct {
	Next     string
	Previous string
	Results  []location_area
}

type location_area struct {
	Name string
	Url  string
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
	}
	return registry
}

func main() {
	regi := generateRegistery()
	c := config{
		prev:  "",
		next:  "https://pokeapi.co/api/v2/location-area/",
		cache: pokecache.NewCache(20 * time.Second),
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		fmt.Println()

		input := cleanInput(scanner.Text())
		inputtedCommand := input[0]

		command, exists := regi[inputtedCommand]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}
		err := command.callback(&c)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println()
	}
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	words := strings.Fields(lowerText)
	return words // first puts everything into lowercase, then
	// seperates every word by whitespaces into a []string
}
