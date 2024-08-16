package main

import (
	"fmt"
	"os"

	"github.com/Thwani47/termilearn/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")

	if err != nil {
		fmt.Println("fatal", err)
		os.Exit(1)
	}

	defer f.Close()

	m := tui.NewMainModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
