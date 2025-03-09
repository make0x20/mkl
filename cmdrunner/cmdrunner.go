package cmdrunner

import (
	"fmt"
	"os"
	"syscall"
)

// RunCmd executes a command
func RunCmd(cmd string) error {
	// Get the user's current shell
	shell, err := getCurrentShell()
	if err != nil {
		return err
	}

	// Execute the command with the shell
	err = syscall.Exec(shell, []string{shell, "-i", "-c", cmd}, os.Environ())
	if err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	return nil
}

// getCurrentShell gets $SHELL environment variable
func getCurrentShell() (string, error) {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell, nil
	}

	return "", fmt.Errorf("SHELL environment variable not set")
}
