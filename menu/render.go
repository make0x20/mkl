package menu

import (
    "github.com/charmbracelet/huh"
    "sort"
)

// menuFormData holds all the prepared menu data
type menuFormData struct {
    Form         *huh.Form
    Selection    *int
    Commands     map[int]string
    Submenus     map[int]map[string]*MenuItem
    SubmenuNames map[int]string
}

// renderMenu displays a menu and returns the selected command
func renderMenu(title string, items map[string]*MenuItem, tc ThemeConfig) (string, error) {
	// Prepare menu form
	menuData, err := prepareMenuForm(title, items, tc)
	if err != nil {
		return "", err
	}
    
    // Run the form
    if err := menuData.Form.Run(); err != nil {
        return "", err
    }
    
    // Return if direct command selected
    if cmd, ok := menuData.Commands[*menuData.Selection]; ok {
        return cmd, nil
    }
    
    // Handle submenu navigation
    if submenu, ok := menuData.Submenus[*menuData.Selection]; ok {
        // Get submenu name for the title
        submenuName := menuData.SubmenuNames[*menuData.Selection]
        // Navigate to submenu
        return renderMenu(submenuName, submenu, tc)
    }
    
    return "", nil
}

// prepareMenuForm creates the menu form and returns menu data
func prepareMenuForm(title string, items map[string]*MenuItem, tc ThemeConfig) (*menuFormData, error) {
	// Create theme
	theme, err := tc.CreateTheme()
	if err != nil {
		return nil, err
	}

    // Track commands and submenus
    commands := make(map[int]string)
    submenus := make(map[int]map[string]*MenuItem)
    submenuNames := make(map[int]string)
    
    // Build options for current menu level
    options := []huh.Option[int]{}
    nextID := 1
    
    // Get all keys and sort them alphabetically for consistent display
    keys := make([]string, 0, len(items))
    for name := range items {
        keys = append(keys, name)
    }

	// Sort keys alphabetically
    sort.Strings(keys)
    
    // Add direct commands
    for _, name := range keys {
        item := items[name]
        if !item.IsSubmenu {
            id := nextID
            commands[id] = item.Command
            options = append(options, huh.NewOption(name, id))
            nextID++
        }
    }
    
    // Add submenus
    for _, name := range keys {
        item := items[name]
        if item.IsSubmenu {
            id := nextID
            submenus[id] = item.Children
            submenuNames[id] = name
            options = append(options, huh.NewOption(name+tc.SubmenuPointer, id))
            nextID++
        }
    }
    
    // Create select component
    var selection int
    selectionPtr := &selection
    selectInput := huh.NewSelect[int]().
        Title(title).
        Options(options...).
        Value(selectionPtr)
    
    // Create form
    form := huh.NewForm(
        huh.NewGroup(selectInput),
    )
    form.WithTheme(theme)
    
    return &menuFormData{
        Form:         form,
        Selection:    selectionPtr,
        Commands:     commands,
        Submenus:     submenus,
        SubmenuNames: submenuNames,
    } , nil
}
