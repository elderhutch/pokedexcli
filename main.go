package main

import (
	"fmt"

	"github.com/elderhutch/pokedexcli/internal/pokecache"
)

func main() {
	pokecache.NewCache(60) // Initialize the cache with a 5-second interval
	// Start the REPL for the Pokedex
	fmt.Print("Pokedex > ")
	cfg := &mapDirection{}
	startRepl(cfg)
}
