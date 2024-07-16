package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Thwani47/termilearn/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type editorFinishedMsg struct {
	err error
}

func openEditor() tea.Cmd {
	editor := "vim" // this should be configurable (vim, nano)
	file := "helloworld.go"

	c := exec.Command(editor, file)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

func main() {
	m := tui.NewMainModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
