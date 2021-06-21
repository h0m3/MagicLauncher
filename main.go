package main

import (
	"log"
	"time"
)

func main() {
	// game, err := decodeGame("/home/artur/test.json")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%v+", game)
	launch("C:\\test.json")
}

func launch(path string) {

	// Process game JSON
	log.Printf("Decoding JSON")
	game, err := decodeGame(path)
	if err != nil {
		log.Fatalf("Unable to decode JSON '%s', %s.\n", path, err)
	}

	// Start the application
	log.Printf("Starting %s ('%s')", game.Name, game.Appid)
	cmd, err := startApplication(game)
	if err != nil {
		log.Fatalf("Unable to start application '%s', %s\n", game.Appid, err)
	}

	// Check if the application started
	log.Printf("Waiting for '%s' to start", game.Appid)
	for i := game.Timeout.Startup; ; i-- {
		if isRunning(game.Process...) {
			break
		}
		if i <= 0 {
			if err := stopApplication(game, cmd); err != nil {
				log.Printf("Unable to stop process, %s", err)
			}
			log.Fatalf("Unable to start application '%s', %d seconds timeout expired\n", game.Appid, game.Timeout.Startup)
		}
		time.Sleep(1 * time.Second)
	}

	// Wait for application to finish
	log.Printf("'%s' is running", game.Appid)
	for isRunning(game.Process...) {
		time.Sleep(5 * time.Second)
	}

	// Shutdown application
	log.Printf("Shutting down '%s'", game.Appid)
	if err := stopApplication(game, cmd); err != nil {
		log.Fatalf("Unable to stop application '%s', %s", game.Appid, err)
	}
}
