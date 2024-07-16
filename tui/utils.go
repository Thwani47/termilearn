package tui

import tea "github.com/charmbracelet/bubbletea"

func IsQuitting(msg tea.KeyMsg) bool {
	if str := msg.String(); str == "ctrl+c" || str == "q" {
		return true
	}
	return false
}
