package tui

import tea "github.com/charmbracelet/bubbletea"

type conceptsModelstate int

const (
	viewConceptsList conceptsModelstate = iota
	viewConcept
)

type conceptsModel struct {
	conceptsList tea.Model
	concept      singleConceptModel
	state        conceptsModelstate
	width        int
	height       int
}

func (m conceptsModel) Init() tea.Cmd {
	return nil
}

func (m conceptsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case viewConceptsList:
		m.conceptsList, cmd = m.conceptsList.Update(msg)
	case viewConcept:
		var model tea.Model
		model, cmd = m.concept.Update(msg)
		m.concept = model.(singleConceptModel)
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if IsQuitting(msg) {
			cmd = tea.Quit
		}
	case conceptSelectedMessage:
		m.state = viewConcept
		m.concept = NewSingleConceptModel(msg.id, msg.choice)
		var model tea.Model
		model, _ = m.concept.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		m.concept = model.(singleConceptModel)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, cmd
}

func (m conceptsModel) View() string {
	switch m.state {
	case viewConceptsList:
		return m.conceptsList.View()
	case viewConcept:
		return m.concept.View()
	default:
		return "I don't think we should be here"
	}
}

func NewConceptsModel() conceptsModel {
	conceptsList := NewConceptListModel()
	return conceptsModel{
		conceptsList: conceptsList,
		state:        viewConceptsList,
	}
}
