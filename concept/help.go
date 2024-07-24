package concept

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Help key.Binding
	Back key.Binding
	Quit key.Binding
}

// ShortHelp returns key bindings to be shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Back, k.Quit}
}

// FullHelp returns key bindings for the expanded help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Back, k.Quit},
	}
}

var keys = keyMap{
	Help: key.NewBinding(key.WithKeys("?", "h"), key.WithHelp("?/h", "toggle help")),
	Back: key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "go back to the concepts list")),
	Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit the application")),
}
