package main

import (
	"fmt"
	"log"
	"github.com/make0x20/mkl/cmdrunner"
	"github.com/make0x20/mkl/config"
	"github.com/make0x20/mkl/menu"
	"github.com/make0x20/mkl/ui"
)

func main() {
	// Parse command line flags
	cfg := config.ParseFlags()

	// Read and parse menu file
	menuBytes, err := menu.ReadMenuFile(cfg.MenuFile)
	if err != nil {
		log.Fatal(err)
	}

	// Parse menu data
	m, err := menu.ParseMenu(menuBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Navigate menu and get selected command
	cmd, err := ui.NavigateMenu(cfg.MenuTitle, m.Items)
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
