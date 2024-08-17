package keys

import "github.com/charmbracelet/bubbles/key"

type PracticeKeyMap struct {
	Help key.Binding
	Back key.Binding
	Quit key.Binding
	Open key.Binding
	Test key.Binding
}

type QuestionListKeyMap struct {
	Help     key.Binding
	Back     key.Binding
	Quit     key.Binding
	Practice key.Binding
}

func (q QuestionListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{q.Practice, q.Help, q.Back, q.Quit}
}

func (q QuestionListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{q.Practice, q.Help},
		{q.Back, q.Quit},
	}
}

func (k PracticeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Open, k.Test, k.Help, k.Back, k.Quit}
}

func (k PracticeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Open, k.Test},
		{k.Help, k.Back, k.Quit},
	}
}

var PracticeKeys = PracticeKeyMap{
	Help: key.NewBinding(key.WithKeys("?", "h"), key.WithHelp("?/h", "toggle help")),
	Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit the application")),
	Back: key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "go back")),
	Open: key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open the practice file for editing")),
	Test: key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "run the tests")),
}

var QuestionListKeys = QuestionListKeyMap{
	Help:     key.NewBinding(key.WithKeys("?", "h"), key.WithHelp("?/h", "toggle help")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit the application")),
	Back:     key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "go back")),
	Practice: key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "start practicing question")),
}
