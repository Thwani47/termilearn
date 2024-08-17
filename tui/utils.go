package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type BackHandler func(tea.WindowSizeMsg) (tea.Model, tea.Cmd)

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
