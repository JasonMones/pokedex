package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UnmarshalFromPokeapi[T any](s *T, url string) ([]byte, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting resources: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body) //read the data from the response
	if err != nil {
		return nil, fmt.Errorf("error reading json: %w", err)
	}

	if err = json.Unmarshal(data, s); err != nil { //grab needed data
		return nil, fmt.Errorf("command unknown")
	}
	return data, nil
}
