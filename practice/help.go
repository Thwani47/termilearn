package practice

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Help key.Binding
	Back key.Binding
	Quit key.Binding
	Open key.Binding
	Test key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Open, k.Test, k.Help, k.Back, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Open, k.Test},
		{k.Help, k.Back, k.Quit},
	}
}

var keys = keyMap{
	Help: key.NewBinding(key.WithKeys("?", "h"), key.WithHelp("?/h", "toggle help")),
	Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit the application")),
	Back: key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "go back")),
	Open: key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open the practice file for editing")),
	Test: key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "run the tests")),
}
