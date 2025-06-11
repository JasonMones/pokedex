package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

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

	offset := numbers[len(numbers)-2] //finds the offset to be used for the previous page
	offset_int, err := strconv.Atoi(offset)
	if err != nil {
		return fmt.Errorf("error parsing int: %w", err)
	}

	c.next = c.prev
	c.prev = fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%s&limit=20", strconv.Itoa(offset_int-40))
	commandMap(c)
	return nil
}

func commandMap(c *config) error {

	var loc_batch location_batch
	loc_batch.Results = make([]location_area, 20)

	if val, ok := c.cache.Get(c.next); ok { //check if url is in cache
		if err := json.Unmarshal(val, &loc_batch); err != nil { //grab needed data
			return fmt.Errorf("error grabbing data: %w", err)
		}

		for _, loc := range loc_batch.Results {
			fmt.Println(loc.Name)
		}
		fmt.Println("returning a cached result")
		c.prev = loc_batch.Previous
		c.next = loc_batch.Next
		return nil
	}

	res, err := http.Get(c.next) //get data from the current resourse page
	if err != nil {
		return fmt.Errorf("error getting resources: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body) //read the data from the response
	if err != nil {
		return fmt.Errorf("error reading json: %w", err)
	}

	if err = json.Unmarshal(data, &loc_batch); err != nil { //grab needed data
		return fmt.Errorf("error grabbing data: %w", err)
	}

	for _, loc := range loc_batch.Results {
		fmt.Println(loc.Name)
	}
	c.cache.Add(c.next, data)

	c.prev = loc_batch.Previous
	c.next = loc_batch.Next
	return nil
}
