package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

// List of process from Steam platform
var steamProcessList []string = []string{
	"steam.exe",
	"steamwebhelper.exe",
}

// Return the Steam executable path
func getSteamPath() string {
	return "C:\\Steam\\steam.exe"
}

// Execute a Steam app given a AppID
func startSteamApp(appid uint32, steamParams []string) (*exec.Cmd, error) {
	appidString := strconv.Itoa(int(appid))
	params := []string{"-applaunch", appidString}
	params = append(params, steamParams...)
	return startCommand(getSteamPath(), params...)
}

// Exit steam gracefully (if possible)
func exitSteam() error {
	return exec.Command(getSteamPath(), "-shutdown").Run()
}

// Launch an application then finish steam
func launchSteamApp(appid uint32, steamParams []string, processList []string, timeout uint16) error {
	// Launch application
	log.Printf("launching Application ID '%d' on Steam", appid)
	cmd, err := startSteamApp(appid, steamParams)
	if err != nil {
		log.Printf("unable to start AppID '%d' on Steam, error: %s", appid, err)
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("unable to kill Application process, error: %s", err)
		}
		return err
	}

	// Wait to application start
	if !timeoutToStart(timeout, processList) {
		log.Printf("Closing steam if still open")
		if err := exitSteam(); err != nil {
			log.Printf("unable to close Steam, error: %s", err)
		}
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("unable to kill Application process, error: %s", err)
		}
		return fmt.Errorf("unable to start '%d' on Steam, %d seconds timeout reached", appid, timeout)
	}

	// Monitor open application
	log.Printf("appID '%d' launched on Steam", appid)
	waitApplication(processList...)

	// Try to exit steam
	log.Printf("closing steam")
	if err := exitSteam(); err != nil {
		if err := cmd.Process.Kill(); err != nil {
			return err
		}
		return nil
	}

	// Check if steam was closed
	if timeoutToStop(timeout, steamProcessList) {
		return nil
	}

	// Force finish all processes
	log.Printf("Tiemout reached, killing remaining processes")
	return cmd.Process.Kill()
}
