package main

// Structs for parsing the API response
type locationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type mapDirection struct {
	next string
	prev string
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
	direction   *mapDirection
}
