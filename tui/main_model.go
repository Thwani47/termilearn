package tui

import tea "github.com/charmbracelet/bubbletea"

type appState int

const (
	viewConceptList appState = iota
	viewConcept
)

type mainModel struct {
	conceptsList tea.Model
	concept      conceptModel
	state        appState
	width        int
	height       int
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case viewConceptList:
		m.conceptsList, cmd = m.conceptsList.Update(msg)
	case viewConcept:
		var model tea.Model
		model, cmd = m.concept.Update(msg)
		m.concept = model.(conceptModel)
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if IsQuitting(msg) {
			cmd = tea.Quit
		}
	case conceptSelectedMessage:
		m.state = viewConcept
		m.concept = NewConceptModel(msg.id, msg.choice)
		var model tea.Model
		model, _ = m.concept.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		m.concept = model.(conceptModel)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, cmd
}

func (m mainModel) View() string {
	switch m.state {
	case viewConceptList:
		return m.conceptsList.View()
	case viewConcept:
		return m.concept.View()
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
