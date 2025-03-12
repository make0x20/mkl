package config

import (
	"flag"
	"os"
	"path/filepath"
)

// Config holds the config
type Config struct {
	MenuFile  string
	ThemeFile string
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

	// Config file path
	flag.StringVar(&config.ThemeFile, "t",
		filepath.Join(home, ".config", "mookielauncher", "theme.json"),
		"Path to the theme file")

	// Menu file path
	flag.StringVar(&config.MenuFile, "m",
		filepath.Join(home, ".config", "mookielauncher", "menu.json"),
		"Path to the menu file")

	// Menu title
	flag.StringVar(&config.MenuTitle, "n",
		"Run command",
		"Name of the menu")

	// Print only
	flag.BoolVar(&config.PrintOnly, "p", false,
		"Print the command instead of running it")

	flag.Parse()

	return config
}
