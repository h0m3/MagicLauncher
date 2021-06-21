package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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
	case "epic":
		return programFilesx86("Epic Games\\Launcher\\Portal\\Binaries\\Win32\\EpicGamesLauncher.exe")
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
		arguments = []string{"-silent", "-applaunch", game.Appid}
	case "epic":
		command = "rundll32"
		url := fmt.Sprintf("com.epicgames.launcher://apps/%s?action=launch&silent=true", game.Appid)
		arguments = []string{"url.dll,FileProtocolHandler", url}
	default:
		return nil, fmt.Errorf("invalid platform '%s'", game.Platform)
	}

	// Create command object and set default parameters
	arguments = append(arguments, game.StartupArguments...)
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
	// Check if autoclose is enabled
	if !game.Autoclose {
		return cmd.Wait()
	}

	// Execute correct action based on the platform
	switch game.Platform {
	case "steam":
		if err := runCommand(game.LauncherPath, "-shutdown"); err != nil {
			return err
		}
		return cmd.Wait()
	case "epic":
		time.Sleep(time.Duration(game.Timeout.Shutdown) * time.Second)
		if err := killProcess("EpicGamesLauncher.exe"); err != nil {
			return err
		}
		if err := killProcess("EpicWebHelper.exe"); err != nil {
			return err
		}
		if err := cmd.Process.Kill(); err != nil {
			return err
		}
		return cmd.Process.Release()
	}

	return fmt.Errorf("invalid platform '%s'", game.Platform)
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

// Run a command, block execution
func runCommand(command string, arguments ...string) error {
	cmd := exec.Command(command, arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Kill a process based on name
func killProcess(process string) error {
	// Get a list of processes
	processes, err := ps.Processes()
	if err != nil {
		return err
	}

	// Find and kill the process in the table
	for _, item := range processes {
		if item.Executable() == process {
			osProcess, err := os.FindProcess(item.Pid())
			if err != nil {
				return err
			}
			if err := osProcess.Kill(); err != nil {
				return err
			}
			return nil
		}
	}

	// Return error if no process found
	return fmt.Errorf("no process found with name '%s'", process)
}
