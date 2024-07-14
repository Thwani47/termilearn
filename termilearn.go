package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	title     string
	err       error
	altScreen bool
}

type editorFinishedMsg struct {
	err error
}

func InitModel() model {
	return model{
		title: "Termilearn",
	}
}

func (m model) Init() tea.Cmd {
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.altScreen = !m.altScreen
			cmd := tea.EnterAltScreen
			if !m.altScreen {
				cmd = tea.ExitAltScreen
			}

			return m, cmd
		case "e":
			return m, openEditor()
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	return m.title
}

func main() {
	if _, err := tea.NewProgram(InitModel()).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
