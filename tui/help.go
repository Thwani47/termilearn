package tui

import "github.com/charmbracelet/bubbles/key"

type tabsKeyMap struct {
	Help   key.Binding
	Next   key.Binding
	Prev   key.Binding
	Quit   key.Binding
	Select key.Binding
}

func (k tabsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Next, k.Prev, k.Quit, k.Help}
}

func (k tabsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Next, k.Prev},
		{k.Quit, k.Help},
	}
}

var tabKeys = tabsKeyMap{
	Help:   key.NewBinding(key.WithKeys("?", "h"), key.WithHelp("?/h", "toggle help")),
	Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit the application")),
	Next:   key.NewBinding(key.WithKeys("right", "n", "tab"), key.WithHelp("tab/n", "move to the next tab")),
	Prev:   key.NewBinding(key.WithKeys("left", "p", "shift+tab"), key.WithHelp("p/shift+tab", "move to the previous tab")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "selected current tab")),
}

type conceptListKeyMap struct {
	Choose key.Binding
	Back   key.Binding
	Quit   key.Binding
}

var conceptListKeys = conceptListKeyMap{
	Choose: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select concept")),
	Back:   key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "go back")),
}
