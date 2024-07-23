package tui

import (
	"fmt"
	"strings"

	"github.com/Thwani47/termilearn/helpers"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const userHighPerformanceRender = false

var (
	viewPortTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	viewPortInfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type conceptModel struct {
	conceptId string
	title     string
	viewport  viewport.Model
	w         tea.WindowSizeMsg
	back      BackHandler
}

func (m conceptModel) Init() tea.Cmd {
	return nil
}

func (m conceptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if IsQuitting(msg) {
			return m, tea.Quit
		}
		switch keypress := msg.String(); keypress {
		case "b":
			return m.back(tea.WindowSizeMsg{Width: 100, Height: 30})
		}
	case tea.WindowSizeMsg:
		m.w = msg
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func readNotes(conceptId string) string {
	notes, err := helpers.ReaderConceptNotes(conceptId)

	if err != nil {
		return fmt.Sprintf("Error reading notes: %s", err)
	}

	return notes
}

func (m conceptModel) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m conceptModel) headerView() string {
	title := viewPortTitleStyle.Render(m.title)
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Height(title)))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m conceptModel) footerView() string {
	info := viewPortInfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func NewConcept(id, title string, width, height int, backHandler BackHandler) (tea.Model, tea.Cmd) {
	// read notes
	content := readNotes(id)
	headerHeight := lipgloss.Height(viewPortTitleStyle.Render(title))
	footerHeight := lipgloss.Height(viewPortInfoStyle.Render(""))

	verticalMarginHeight := headerHeight + footerHeight
	availableHeight := height - verticalMarginHeight

	vp := viewport.New(width, availableHeight)
	vp.YPosition = headerHeight
	vp.HighPerformanceRendering = userHighPerformanceRender
	vp.SetContent(readNotes(content))

	m := conceptModel{
		conceptId: id,
		title:     title,
		w: tea.WindowSizeMsg{
			Width:  width,
			Height: height,
		},
		viewport: vp,
		back:     backHandler,
	}

	cmd := tea.SetWindowTitle(title)
	return m, cmd

}
