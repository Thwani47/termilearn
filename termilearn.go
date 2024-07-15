package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

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

type editorFinishedMsg struct {
	err error
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func openEditor() tea.Cmd {
	editor := "vim" // this should be configurable (vim, nano)
	file := "helloworld.go"

	c := exec.Command(editor, file)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
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
	}
	return m, cmd
}

func (m mainModel) View() string {
	switch m.state {
	case viewConceptList:
		return m.conceptsList.View()
	default:
		return "Hello world"
	}
}

func NewMainModel() mainModel {
	conceptsList := NewConceptListModel()
	return mainModel{
		conceptsList: conceptsList,
		state:        viewConceptList,
	}
}

func main() {
	m := NewMainModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
