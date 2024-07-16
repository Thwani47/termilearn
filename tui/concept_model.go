package tui

import tea "github.com/charmbracelet/bubbletea"

type conceptModel struct {
	Title       string
	Description string
}

func (m conceptModel) Init() tea.Cmd {
	return nil
}

func (m conceptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m conceptModel) View() string {
	return m.Title + "\n\n" + m.Description
}

func NewConceptModel(title, description string) conceptModel {
	return conceptModel{
		Title:       title,
		Description: description,
	}
}
