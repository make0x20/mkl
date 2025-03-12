package main

import (
	"errors"
	"fmt"
	"github.com/make0x20/mkl/cmdrunner"
	"github.com/make0x20/mkl/config"
	"github.com/make0x20/mkl/menu"
	"log"
	"os"
)

func main() {
	// Parse command line flags
	cfg := config.ParseFlags()

	// Read and parse menu file
	menuBytes, err := menu.ReadMenuFile(cfg.MenuFile)
	if err != nil {
		log.Fatal(err)
	}

	// Read and parse theme file
	themeBytes, err := menu.ReadThemeFile(cfg.ThemeFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	// Parse menu data
	m, err := menu.NewMenu(menuBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Parse theme data
	theme, err := menu.NewThemeConfig(string(themeBytes))
	if err != nil {
		log.Fatal(err)
	}

	// Navigate menu and get selected command
	cmd, err := m.Render(cfg.MenuTitle, theme)
	if err != nil {
		log.Fatal(err)
	}

	// Handle the selected command
	if cmd == "" {
		return // No command selected
	}

	if cfg.PrintOnly {
		fmt.Println(cmd)
		return
	}

	// Execute the command
	if err := cmdrunner.RunCmd(cmd); err != nil {
		log.Fatal(err)
	}
}
