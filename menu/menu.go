package menu

import (
	"encoding/json"
	"fmt"
	"os"
)

// MenuItem holds all the menu item data
type MenuItem struct {
	Name      string               // Name of the menu item
	Command   string               // Command to run
	IsSubmenu bool                 // Whether the item is a submenu
	Children  map[string]*MenuItem // Submenu items
}

// Menu holds a map of menu items
type Menu struct {
	Items map[string]*MenuItem
}

// ReadMenuFile reads a menu file and returns bytes
func ReadMenuFile(file string) ([]byte, error) {
	menuData, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading menu file: %v", err)
	}
	return menuData, nil
}

// ParseMenu parses a menu byte slice into a Menu struct
func ParseMenu(data []byte) (*Menu, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty menu data")
	}

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return nil, fmt.Errorf("failed to parse menu JSON: %w", err)
	}

	menu := &Menu{
		Items: make(map[string]*MenuItem),
	}

	for name, rawValue := range rawMap {
		item, err := parseMenuItem(name, rawValue)
		if err != nil {
			return nil, err
		}
		menu.Items[name] = item
	}

	return menu, nil
}

// parseMenuItem parses a raw JSON message into a MenuItem
func parseMenuItem(name string, data json.RawMessage) (*MenuItem, error) {
	// Try to unmarshal as string first (command)
	var commandStr string
	if err := json.Unmarshal(data, &commandStr); err == nil {
		return &MenuItem{
			Name:      name,
			Command:   commandStr,
			IsSubmenu: false,
		}, nil
	}

	// If not a string, try as map (submenu)
	var submenuMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &submenuMap); err != nil {
		return nil, fmt.Errorf("could not parse menu item %q: %w", name, err)
	}

	submenu := &MenuItem{
		Name:      name,
		IsSubmenu: true,
		Children:  make(map[string]*MenuItem),
	}

	for childName, childData := range submenuMap {
		child, err := parseMenuItem(childName, childData)
		if err != nil {
			return nil, err
		}
		submenu.Children[childName] = child
	}

	return submenu, nil
}

