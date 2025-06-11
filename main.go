package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
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

type location_area struct {
	Name string
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
		next:  "https://pokeapi.co/api/v2/location-area/1/",
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

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	regi := generateRegistery()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range regi {
		fmt.Printf("%s: %s", command.name, command.description)
		fmt.Println()
	}
	return nil
}

func commandMapb(c *config) error {
	if c.prev == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	re := regexp.MustCompile("[0-9]+") //find digits
	numbers := re.FindAllString(c.next, -1)

	URLresource, err := strconv.Atoi(numbers[len(numbers)-1]) //turns the last digit found into an int for use
	if err != nil {
		return fmt.Errorf("error parsing int: %w", err)
	}

	if URLresource < 40 {
		c.prev = ""
		c.next = "https://pokeapi.co/api/v2/location-area/1/"
		fmt.Println("you're on the first page")
		return nil
	}
	c.next = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", strconv.Itoa(URLresource-40))
	commandMap(c)
	return nil
}

func commandMap(c *config) error {

	re := regexp.MustCompile("[0-9]+") //find digits
	numbers := re.FindAllString(c.next, -1)

	URLresourceIndex, err := strconv.Atoi(numbers[len(numbers)-1]) //turns the last digit found into an int for use
	if err != nil {
		return fmt.Errorf("error parsing int: %w", err)
	}

	baseURL := "https://pokeapi.co/api/v2/location-area"

	successfulResourcePulls := 0
	var loc location_area
	for successfulResourcePulls <= 20 { //print 20 locations starting from c.next

		currentURL := fmt.Sprintf("%s/%d/", baseURL, URLresourceIndex)

		if val, ok := c.cache.Get(currentURL); ok { //check if url is in cache
			fmt.Println(string(val))
			URLresourceIndex++
			continue
		}

		fmt.Printf("Getting resources from %s...", currentURL)
		res, err := http.Get(currentURL) //get data from the current resourse page
		if err != nil {
			return fmt.Errorf("error getting resources: %w", err)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body) //read the data from the response
		if err != nil {
			return fmt.Errorf("error reading json: %w", err)
		}

		if err = json.Unmarshal(data, &loc); err != nil { //grab needed data
			//return fmt.Errorf("error grabbing data: %w", err)
			fmt.Println("Failed to get resources")
			URLresourceIndex++
			continue
		}
		fmt.Println(loc.Name)
		c.cache.Add(currentURL, []byte(loc.Name)) //add location name to cache
		res.Body.Close()
		URLresourceIndex++
		successfulResourcePulls++
	}

	c.prev = c.next //set config for next call to map
	c.next = fmt.Sprintf("%s/%d/", baseURL, URLresourceIndex)
	return nil
}
