package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BackHandler func(tea.WindowSizeMsg) (tea.Model, tea.Cmd)

func IsQuitting(msg tea.KeyMsg) bool {
	if str := msg.String(); str == "ctrl+c" || str == "q" {
		return true
	}
	return false
}

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

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
