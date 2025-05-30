package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var lastWords []string

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

func captureCommand(s string) error {
	url := "https://pokeapi.co/api/v2/pokemon/" + s
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
	var apiResp ExploredArea
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&apiResp)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", apiResp.Name)

	return nil
}

func exploreCommand(s string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + s
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
	var apiResp ExploredArea
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&apiResp)
	if err != nil {
		return err
	}

	fmt.Printf("Location Area: %s\n", apiResp.Name)
	fmt.Println("Pokemon encounters:")
	for _, encounter := range apiResp.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
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
			fmt.Println("You're on the first page.")
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

func startRepl(cfg *mapDirection) {
	scanner := bufio.NewScanner(os.Stdin)
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
		"explore": {
			name:        "explore",
			description: "Explore a location",
			callback: func() error {
				if len(lastWords) > 1 {
					location := strings.TrimSpace(lastWords[1])
					return exploreCommand(location)
				}
				fmt.Print("Enter location area name: ")
				inputScanner := bufio.NewScanner(os.Stdin)
				if inputScanner.Scan() {
					location := strings.TrimSpace(inputScanner.Text())
					return exploreCommand(location)
				}
				return nil
			},
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback: func() error {
				if len(lastWords) > 1 {
					pokemon := strings.TrimSpace(lastWords[1])
					return captureCommand(pokemon)
				}
				fmt.Print("Enter a pokemon name: ")
				inputScanner := bufio.NewScanner(os.Stdin)
				if inputScanner.Scan() {
					pokemon := strings.TrimSpace(inputScanner.Text())
					return captureCommand(pokemon)
				}
				return nil
			},
			direction: cfg,
		},
	}

	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		lastWords = words
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
