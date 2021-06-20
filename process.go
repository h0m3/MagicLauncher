package main

import (
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/process"
)

// Start a command a return the cmd control
// Dont block execution
func startCommand(program string, params ...string) (*exec.Cmd, error) {
	cmd := exec.Command(program, params...)
	err := cmd.Start()
	return cmd, err
}

// Open a application URL using startCommand()
// Return cmd control
func openURL(url string) (*exec.Cmd, error) {
	return startCommand("rundll32", "url.dll,FileProtocolHandler", url)
}

// Check if a list of processes is running
// Return true if any process is running, false otherwise
func isRunning(names ...string) bool {
	processes, err := process.Processes()
	if err != nil {
		return false
	}

	for _, process := range processes {
		n, err := process.Name()
		if err != nil {
			continue
		}

		for _, name := range names {
			if n == name {
				return true
			}
		}
	}

	return false
}

// Check if a list of process started until the timeout ends
// Return false if the timeout reached end and true if found the process
func timeoutToStart(timeout uint16, programs []string) bool {
	for i := uint16(0); i < timeout; i++ {
		if isRunning(programs...) {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

// Check if a list of process ended until the timeout ends
// Return false if the timeout reached end and true if found process
func timeoutToStop(timeout uint16, programs []string) bool {
	for i := uint16(0); i < timeout; i++ {
		if !isRunning(programs...) {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

// Kill a list of processes in order
// If a process dont exists (already killed) do nothing
// func killProcess(names ...string) error {
// 	processes, err := process.Processes()
// 	if err != nil {
// 		return err
// 	}

// 	for _, process := range processes {

// 		// Check if the process still exists and its running
// 		running, err := process.IsRunning()
// 		if err != nil || !running {
// 			continue
// 		}

// 		n, err := process.Name()
// 		if err != nil {
// 			return err
// 		}
// 		for _, name := range names {
// 			if n == name {
// 				if err := process.Kill(); err != nil {
// 					return err
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

// Wait for application
func waitApplication(processes ...string) {
	for isRunning(processes...) {
		time.Sleep(5 * time.Second)
	}
}

// func getProcess(names ...string) ([]*process.Process, error) {
// 	var foundProcesses []*process.Process

// 	processes, err := process.Processes()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, process := range processes {
// 		n, err := process.Name()
// 		if err != nil {
// 			return nil, err
// 		}

// 		for _, name := range names {
// 			if n == name {
// 				foundProcesses = append(foundProcesses, process)
// 			}
// 		}
// 	}

// 	if foundProcesses == nil {
// 		return nil, fmt.Errorf("no process found")
// 	}

// 	return foundProcesses, nil
// }
