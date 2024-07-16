package tui

import tea "github.com/charmbracelet/bubbletea"

type appState int

const (
	viewConceptList appState = iota
	viewConcept
)

type mainModel struct {
	conceptsList tea.Model
	state        appState
	title        string
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case viewConceptList:
		m.conceptsList, cmd = m.conceptsList.Update(msg)
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmd = tea.Quit
		}
	case conceptSelectedMessage:
		m.title = msg.choice
		m.state = viewConcept
	}
	return m, cmd
}

func (m mainModel) View() string {
	switch m.state {
	case viewConceptList:
		return m.conceptsList.View()
	default:
		return "I don't think we should be here"
	}
}

func NewMainModel() mainModel {
	conceptsList := NewConceptListModel()
	return mainModel{
		conceptsList: conceptsList,
		state:        viewConceptList,
	}
}
