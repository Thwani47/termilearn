package main

import (
	"fmt"
	"github.com/Thwani47/termilearn/tui"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	m := tui.NewMainModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
