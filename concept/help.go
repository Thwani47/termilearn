package concept

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Help     key.Binding
	Up       key.Binding
	Down     key.Binding
	Back     key.Binding
	Quit     key.Binding
	Practice key.Binding
}

// ShortHelp returns key bindings to be shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Practice, k.Up, k.Down, k.Help, k.Back, k.Quit}
}

// FullHelp returns key bindings for the expanded help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Practice, k.Up, k.Down},
		{k.Help, k.Back, k.Quit},
	}
}

var keys = keyMap{
	Practice: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "practice")),
	Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "scroll up")),
	Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "scroll down")),
	Help:     key.NewBinding(key.WithKeys("?", "h"), key.WithHelp("?/h", "toggle help")),
	Back:     key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "go back to the concepts list")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit the application")),
}
