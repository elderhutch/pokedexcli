package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

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

func helpCommand() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

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

func mapCommand(cfg *mapDirection) func() error {
	return func() error {
		url := "https://pokeapi.co/api/v2/location-area"
		if cfg.next != "" {
			url = cfg.next
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		var apiResp locationAreaResponse
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&apiResp)
		if err != nil {
			return err
		}

		cfg.next = apiResp.Next
		if apiResp.Previous != nil {
			if prevStr, ok := apiResp.Previous.(string); ok {
				cfg.prev = prevStr
			} else {
				cfg.prev = ""
			}
		} else {
			cfg.prev = ""
		}

		for _, area := range apiResp.Results {
			fmt.Println(area.Name)
		}
		return nil
	}
}

func mapbCommand(cfg *mapDirection) func() error {
	return func() error {
		if cfg.prev == "" {
			fmt.Println("No previous page available.")
			return nil
		}
		req, err := http.NewRequest("GET", cfg.prev, nil)
		if err != nil {
			return err
		}
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		var apiResp locationAreaResponse
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&apiResp)
		if err != nil {
			return err
		}

		cfg.next = apiResp.Next
		if apiResp.Previous != nil {
			if prevStr, ok := apiResp.Previous.(string); ok {
				cfg.prev = prevStr
			} else {
				cfg.prev = ""
			}
		} else {
			cfg.prev = ""
		}

		for _, area := range apiResp.Results {
			fmt.Println(area.Name)
		}
		return nil
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")
	cfg := &mapDirection{}
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Show help",
			callback:    helpCommand,
		},
		"map": {
			name:        "map",
			description: "Show map",
			callback:    mapCommand(cfg),
			direction:   cfg,
		},
		"mapb": {
			name:        "mapb",
			description: "Show previous map",
			callback:    mapbCommand(cfg),
			direction:   cfg,
		},
	}

	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		if len(words) == 0 {
			fmt.Print("Pokedex > ")
			continue
		}
		cmdName := strings.ToLower(words[0])
		cmd, exists := commands[cmdName]
		if exists {
			cmd.callback()
		} else {
			fmt.Println("Unknown command:", cmdName)
		}
		fmt.Print("Pokedex > ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
