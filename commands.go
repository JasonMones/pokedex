package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type location_batch struct {
	Next     string
	Previous string
	Results  []location_area
}

type location_area struct {
	Name string
	Url  string
}

type pokemon_batch struct {
	Pokemon_encounters []pokemon_encounters
}

type pokemon_encounters struct {
	Pokemon pokemon
}

type pokemon struct {
	Name string
	Url  string
}

func commandExit(c *config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("unspecified leading arguments")
	}
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("unspecified leading arguments")
	}

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

func commandMapb(c *config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("unspecified leading arguments")
	}

	if c.prev == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	re := regexp.MustCompile("[0-9]+") //find digits
	numbers := re.FindAllString(c.next, -1)

	offset := numbers[len(numbers)-2] //finds the offset to be used for the previous page
	offset_int, err := strconv.Atoi(offset)
	if err != nil {
		return fmt.Errorf("error parsing int: %w", err)
	}

	c.next = c.prev
	c.prev = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%s&limit=20", strconv.Itoa(offset_int-40))
	commandMap(c, args)
	return nil
}

func commandMap(c *config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("unspecified leading arguments")
	}

	var loc_batch location_batch
	loc_batch.Results = make([]location_area, 20)

	if val, ok := c.cache.Get(c.next); ok { //check if url is in cache
		if err := json.Unmarshal(val, &loc_batch); err != nil { //grab needed data
			return fmt.Errorf("error grabbing data: %w", err)
		}

		for _, loc := range loc_batch.Results {
			fmt.Println(loc.Name)
		}
		c.prev = loc_batch.Previous
		c.next = loc_batch.Next
		return nil
	}

	data, err := UnmarshalFromPokeapi(&loc_batch, c.next)
	if err != nil {
		return err
	}

	for _, loc := range loc_batch.Results {
		fmt.Println(loc.Name)
	}

	c.cache.Add(c.next, data)

	c.prev = loc_batch.Previous
	c.next = loc_batch.Next
	return nil
}

func commandExploreHelper(location string, poke_batch *pokemon_batch) {
	fmt.Printf("Exploring %s.", location)
	time.Sleep(500 * time.Millisecond)
	fmt.Print(".")
	time.Sleep(500 * time.Millisecond)
	fmt.Print(".\n")
	fmt.Println("Found Pokemon:")

	for _, pokemon := range poke_batch.Pokemon_encounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
}

func commandExplore(c *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("required arguments not present")
	} else if len(args) > 1 {
		return fmt.Errorf("unspecified leading arguments")
	}
	location := args[0]
	fmt.Println(location)

	var poke_batch pokemon_batch
	URL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location)

	if data, ok := c.cache.Get(URL); ok { //check if url is in cache
		if err := json.Unmarshal(data, &poke_batch); err != nil { //grab needed data
			return fmt.Errorf("error grabbing data: %w", err)
		}

		commandExploreHelper(location, &poke_batch)
		return nil
	}

	data, err := UnmarshalFromPokeapi(&poke_batch, URL)
	if err != nil {
		return err
	}

	c.cache.Add(location, data)

	commandExploreHelper(location, &poke_batch)
	return nil
}
