package menu

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Theme holds all the styles for the UI
type ThemeConfig struct {
	BaseForeground           string `json:"base_foreground,omitempty"`
	BaseBackground           string `json:"base_background,omitempty"`
	BaseBold                 bool   `json:"base_bold,omitempty"`
	TitleForeground          string `json:"title_foreground,omitempty"`
	TitleBackground          string `json:"title_background,omitempty"`
	TitleBold                bool   `json:"title_bold,omitempty"`
	TitleBorder              string `json:"title_separator,omitempty"`
	OptionForeground         string `json:"option_foreground,omitempty"`
	OptionBackground         string `json:"option_background,omitempty"`
	OptionBold               bool   `json:"option_bold,omitempty"`
	SelectedOptionForeground string `json:"selected_foreground,omitempty"`
	SelectedOptionBackground string `json:"selected_background,omitempty"`
	SelectedOptionBold       bool   `json:"selected_bold,omitempty"`
	SelectSelectorForeground string `json:"selector_foreground,omitempty"`
	SelectSelectorBackground string `json:"selector_background,omitempty"`
	SelectSelectorBold       bool   `json:"selector_bold,omitempty"`
	SelectSelectorString     string `json:"selector_string,omitempty"`
	SubmenuPointer           string `json:"submenu_pointer,omitempty"`
}

// Border types
var borderTypes = map[string]lipgloss.Border{
	"normal":      lipgloss.NormalBorder(),
	"thick":       lipgloss.ThickBorder(),
	"double":      lipgloss.DoubleBorder(),
	"block":       lipgloss.BlockBorder(),
	"block_outer": lipgloss.OuterHalfBlockBorder(),
	"block_inner": lipgloss.InnerHalfBlockBorder(),
}

// DefaultTheme returns the default theme
func DefaultThemeConfig() ThemeConfig {
	return ThemeConfig{
		SelectSelectorString: "> ",
		SubmenuPointer:       " >",
	}
}

func NewThemeConfig(themeJSON string) (ThemeConfig, error) {
	theme := DefaultThemeConfig()

	if themeJSON != "" {
		err := json.Unmarshal([]byte(themeJSON), &theme)
		if err != nil {
			return ThemeConfig{}, fmt.Errorf("error unmarshalling custom theme JSON: %w", err)
		}
	}

	return theme, nil
}

// NewTheme takes a custom JSON string and returns a new theme with the customizations applied
func (tc *ThemeConfig) CreateTheme() (*huh.Theme, error) {
	t := huh.ThemeBase()

	// Base style
	t.Focused.Base = lipgloss.NewStyle().
		Foreground(lipgloss.Color(tc.BaseForeground)).
		Bold(tc.BaseBold).
		Padding(1, 1, 0, 1)
	if tc.BaseBackground != "" {
		t.Focused.Base = t.Focused.Base.Background(lipgloss.Color(tc.BaseBackground))
	}

	// Title style
	t.Focused.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(tc.TitleForeground)).
		Bold(tc.TitleBold)
	if tc.TitleBackground != "" {
		t.Focused.Title = t.Focused.Title.Background(lipgloss.Color(tc.TitleBackground))
	}
	// Title border - bottom
	if tc.TitleBorder != "" {
		t.Focused.Title.BorderBottom(false)
	}
	t.Focused.Title = t.Focused.Title.Border(borderTypes[tc.TitleBorder], false, false, true, false)

	// Option style
	t.Focused.UnselectedOption = lipgloss.NewStyle().
		Foreground(lipgloss.Color(tc.OptionForeground)).
		Bold(tc.OptionBold)
	if tc.OptionBackground != "" {
		t.Focused.UnselectedOption = t.Focused.Option.Background(lipgloss.Color(tc.OptionBackground))
	}

	// SelectedOption style
	t.Focused.SelectedOption = lipgloss.NewStyle().
		Foreground(lipgloss.Color(tc.SelectedOptionForeground)).
		Bold(tc.SelectedOptionBold).
		Padding(0, 1)
	if tc.SelectedOptionBackground != "" {
		t.Focused.SelectedOption = t.Focused.SelectedOption.Background(lipgloss.Color(tc.SelectedOptionBackground))
	}

	// SelectSelector style
	t.Focused.SelectSelector = lipgloss.NewStyle().
		Foreground(lipgloss.Color(tc.SelectSelectorForeground)).
		Bold(tc.SelectSelectorBold)
	if tc.SelectSelectorBackground != "" {
		t.Focused.SelectSelector = t.Focused.SelectSelector.Background(lipgloss.Color(tc.SelectSelectorBackground))
	}
	if tc.SelectSelectorString != "" {
		t.Focused.SelectSelector = t.Focused.SelectSelector.SetString(tc.SelectSelectorString)
	}

	return t, nil
}

func ReadThemeFile(file string) ([]byte, error) {
	themeData, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading theme file: %w", err)
	}

	return themeData, nil
}
