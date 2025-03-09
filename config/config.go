package config

import (
	"flag"
	"os"
	"path/filepath"
)

// Config holds the config
type Config struct {
	MenuFile  string
	MenuTitle string
	PrintOnly bool
}

// ParseFlags parses flags into config
func ParseFlags() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	config := Config{}

	flag.StringVar(&config.MenuFile, "m",
		filepath.Join(home, ".config", "mookielauncher", "menu.json"),
		"Path to the menu file")

	flag.StringVar(&config.MenuTitle, "n",
		"Run command - Mookielauncher",
		"Name of the menu")

	flag.BoolVar(&config.PrintOnly, "p", false,
		"Print the command instead of running it")

	flag.Parse()

	return config
}
