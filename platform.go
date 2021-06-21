package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mitchellh/go-ps"
)

// Return path of Program Files (x86)
func programFilesx86(path string) string {
	if programFiles := os.Getenv("PROGRAMFILES(X86"); programFiles != "" {
		return filepath.Join(programFiles, path)
	}
	return filepath.Join("C:\\Program Files (x86)", path)
}

// Return path of Program Files
func programFiles(path string) string {
	if programFiles := os.Getenv("PROGRAMFILES"); programFiles != "" {
		return filepath.Join(programFiles, path)
	}
	return filepath.Join("C:\\Program Files", path)
}

// Return path of a specific platform
func getPlatformPath(platform string) string {
	switch platform {
	case "steam":
		return programFilesx86("Steam\\steam.exe")
	}
	return ""
}

// Run game on specific platform
func startApplication(game Game) (*exec.Cmd, error) {

	// Get correct syntax accourdly with platform
	var command string
	var arguments []string

	switch game.Platform {
	case "steam":
		command = game.LauncherPath
		arguments = []string{"-applaunch", game.Appid}
		arguments = append(arguments, game.StartupArguments...)
	}

	// Create command object and set default parameters
	cmd := exec.Command(command, arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start application and return command object
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

// Stop application and exit platform if necessary
func stopApplication(game Game, cmd *exec.Cmd) error {

	// Check if process is already closed
	if process, err := ps.FindProcess(cmd.Process.Pid); process == nil {
		if err := cmd.Process.Release(); err != nil {
			return err
		}
		return err
	}

	// Get correct syntax accourdly with platform
	var command string
	var arguments []string

	switch game.Platform {
	case "steam":
		command = game.LauncherPath
		arguments = []string{"-shutdown"}
	}

	// Execute stop action
	cmdStop := exec.Command(command, arguments...)
	cmdStop.Stdout = os.Stdout
	cmdStop.Stderr = os.Stderr

	if err := cmdStop.Run(); err != nil {
		return err
	}

	// Wait for process to end
	return cmd.Wait()
}

// Check if a list of processes is running
func isRunning(processes ...string) bool {
	processTable, err := ps.Processes()
	if err != nil {
		return false
	}

	for _, process := range processTable {
		for _, name := range processes {
			if name == process.Executable() {
				return true
			}
		}
	}

	return false
}
