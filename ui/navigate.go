package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/make0x20/mkl/menu"
	"sort"
)

// NavigateMenu displays a menu and returns the selected command
func NavigateMenu(title string, items map[string]*menu.MenuItem) (string, error) {
	// Track commands and submenus
	commands := make(map[int]string)
	submenus := make(map[int]map[string]*menu.MenuItem)
	submenuNames := make(map[int]string)

	// Build options for current menu level
	options := []huh.Option[int]{}
	nextID := 1

	// Get all keys and sort them alphabetically for consistent display
	keys := make([]string, 0, len(items))
	for name := range items {
		keys = append(keys, name)
	}
	sort.Strings(keys) // Sort keys alphabetically

	// Add direct commands (sorted by name)
	for _, name := range keys {
		item := items[name]
		if !item.IsSubmenu {
			id := nextID
			commands[id] = item.Command
			options = append(options, huh.NewOption(name, id))
			nextID++
		}
	}

	// Add submenus (sorted by name)
	for _, name := range keys {
		item := items[name]
		if item.IsSubmenu {
			id := nextID
			submenus[id] = item.Children
			submenuNames[id] = name
			options = append(options, huh.NewOption(name+" >", id))
			nextID++
		}
	}

	// Create select component with custom key bindings
	var selection int
	selectInput := huh.NewSelect[int]().
		Title(title).
		Options(options...).
		Value(&selection)

	// Create and run the form
	form := huh.NewForm(
		huh.NewGroup(selectInput),
	)
	form.WithTheme(GetMenuTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	// Handle selection
	if cmd, ok := commands[selection]; ok {
		// Command selected
		return cmd, nil
	} else if submenu, ok := submenus[selection]; ok {
		// Get submenu name for the title
		submenuName := submenuNames[selection]
		// Navigate to submenu
		return NavigateMenu(submenuName, submenu)
	}

	return "", nil
}
