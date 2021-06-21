package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GameTimeout struct {
	Startup  uint16
	Shutdown uint16
}

type Game struct {
	Appid            string      // Required
	Platform         string      // Required
	Process          []string    // Required
	Timeout          GameTimeout // Optional (default)
	LauncherPath     string      // Optional (default)
	Name             string      // Optional
	StartupArguments []string    // Optional
	Autoclose        bool        // Optional
}

// Create new empty game structure
func newGame() Game {
	return Game{
		Timeout:   GameTimeout{Startup: 60, Shutdown: 60},
		Autoclose: true,
	}
}

// Decode game data based on JSON
func decodeGame(jsonPath string) (Game, error) {
	game := newGame()

	// Read JSON file
	file, err := os.Open(jsonPath)
	if err != nil {
		return game, fmt.Errorf("unable to read game JSON, %s", err)
	}

	// Decode JSON content
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&game)
	if err != nil {
		return game, fmt.Errorf("unable to decode game JSON, %s", err)
	}

	// Check for missing information
	switch "" {
	case game.Appid:
		return game, fmt.Errorf("empty AppID")
	case game.Platform:
		return game, fmt.Errorf("empty platform name")
	}

	if len(game.Process) < 1 {
		return game, fmt.Errorf("no game processes suplied")
	}

	// Add default launcher path if necessary
	if game.LauncherPath == "" {
		game.LauncherPath = getPlatformPath(game.Platform)
	}

	return game, nil
}
