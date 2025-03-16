package menu

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"os"
	"os/exec"
	"runtime"
	"sort"
)

// menuFormData holds all the data needed to render a menu
type menuFormData struct {
	Form         *huh.Form
	Selection    *int
	Commands     map[int]string
	Submenus     map[int]map[string]*MenuItem
	SubmenuNames map[int]string
	GoBack       bool
	Quit         bool
}

// menuModel is a BubbleTea model that wraps our form
type menuModel struct {
	menuData *menuFormData
	form     *huh.Form
}

// newMenuModel creates a new menu model
func newMenuModel(data *menuFormData) menuModel {
	return menuModel{
		menuData: data,
		form:     data.Form,
	}
}

// Init initializes the model
func (m menuModel) Init() tea.Cmd {
	return m.form.Init()
}

// Update handles events and updates the model
func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle key presses
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "h", "left", "esc":
			m.menuData.GoBack = true
			return m, tea.Quit
		case "q", "ctrl+c":
			m.menuData.Quit = true
			return m, tea.Quit
		}
	}

	// Pass messages to the form and handle the returned model
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// Check if the form is completed
	if m.form.State == huh.StateCompleted {
		return m, tea.Quit
	}

	return m, cmd
}

// View renders the model
func (m menuModel) View() string {
	return m.form.View()
}

// renderMenu displays a menu and returns the selected command
func renderMenu(title string, items map[string]*MenuItem, tc ThemeConfig) (string, error) {
	clearScreen()

	menuData, err := prepareMenuForm(title, items, tc)
	if err != nil {
		return "", err
	}

	// Run the BubbleTea program
	p := tea.NewProgram(newMenuModel(menuData))
	result, err := p.Run()
	if err != nil {
		return "", err
	}

	finalModel := result.(menuModel)

	// Handle navigation actions
	if finalModel.menuData.Quit {
		os.Exit(0)
	}

	if finalModel.menuData.GoBack {
		return "", nil // Signal to go back to parent menu
	}

	selection := *finalModel.menuData.Selection

	// Handle command selection
	if cmd, ok := finalModel.menuData.Commands[selection]; ok {
		return cmd, nil
	}

	// Handle submenu navigation
	if submenu, ok := finalModel.menuData.Submenus[selection]; ok {
		submenuName := finalModel.menuData.SubmenuNames[selection]
		cmd, err := renderMenu(submenuName, submenu, tc)

		// Return to this menu if coming back from submenu
		if cmd == "" && err == nil {
			return renderMenu(title, items, tc)
		}

		return cmd, err
	}

	return "", nil
}

// clearScreen clears the terminal screen based on OS
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

// prepareMenuForm creates the menu form with all options
func prepareMenuForm(title string, items map[string]*MenuItem, tc ThemeConfig) (*menuFormData, error) {
	theme, err := tc.CreateTheme()
	if err != nil {
		return nil, err
	}

	menuData := &menuFormData{
		Commands:     make(map[int]string),
		Submenus:     make(map[int]map[string]*MenuItem),
		SubmenuNames: make(map[int]string),
	}

	// Sort menu items for consistent display
	keys := getSortedKeys(items)
	options := buildMenuOptions(keys, items, tc, menuData)

	// Create the form with selection component
	var selection int
	menuData.Selection = &selection

	form := createFormWithKeyBindings(title, options, menuData.Selection, theme)
	menuData.Form = form

	return menuData, nil
}

// getSortedKeys returns alphabetically sorted keys from a map
func getSortedKeys(items map[string]*MenuItem) []string {
	keys := make([]string, 0, len(items))
	for name := range items {
		keys = append(keys, name)
	}
	sort.Strings(keys)
	return keys
}

// buildMenuOptions creates options for the menu form
func buildMenuOptions(keys []string, items map[string]*MenuItem,
	tc ThemeConfig, menuData *menuFormData) []huh.Option[int] {
	options := []huh.Option[int]{}
	nextID := 1

	// Add direct commands first
	for _, name := range keys {
		item := items[name]
		if !item.IsSubmenu {
			menuData.Commands[nextID] = item.Command
			options = append(options, huh.NewOption(name, nextID))
			nextID++
		}
	}

	// add submenus
	for _, name := range keys {
		item := items[name]
		if item.IsSubmenu {
			menuData.Submenus[nextID] = item.Children
			menuData.SubmenuNames[nextID] = name
			options = append(options, huh.NewOption(name+tc.SubmenuPointer, nextID))
			nextID++
		}
	}

	return options
}

// createFormWithKeyBindings creates a form with custom key bindings
func createFormWithKeyBindings(title string, options []huh.Option[int],
	selection *int, theme *huh.Theme) *huh.Form {
	// Create select component
	selectInput := huh.NewSelect[int]().
		Title(title).
		Options(options...).
		Value(selection)

	// Create form with the select component
	form := huh.NewForm(huh.NewGroup(selectInput))

	// Add custom key bindings
	keyMap := &huh.KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("ctrl+c/q", "quit"),
		),
		Select: huh.SelectKeyMap{
			Submit: key.NewBinding(
				key.WithKeys("enter", "right"),
				key.WithHelp("enter/right", "submit"),
			),
			Up: key.NewBinding(
				key.WithKeys("up", "k", "ctrl+k", "ctrl+p"),
				key.WithHelp("↑/k", "up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "j", "ctrl+j", "ctrl+n"),
				key.WithHelp("↓/j", "down"),
			),
			Filter: key.NewBinding(
				key.WithKeys("/"),
				key.WithHelp("/", "filter"),
			),
		},
	}

	return form.WithTheme(theme).WithKeyMap(keyMap)
}

