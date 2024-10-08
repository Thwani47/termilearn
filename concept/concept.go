package concept

import (
	"fmt"
	"strings"

	"github.com/Thwani47/termilearn/common/styles"
	"github.com/Thwani47/termilearn/practice"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BackHandler func(tea.WindowSizeMsg) (tea.Model, tea.Cmd)

const userHighPerformanceRender = false

var helpHeight = 6

type conceptModel struct {
	help      help.Model
	back      BackHandler
	conceptId string
	title     string
	viewport  viewport.Model
	keys      viewportKeyMap
	w         tea.WindowSizeMsg
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
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			return m.back(tea.WindowSizeMsg{Width: 100, Height: 30})
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Practice):
			return practice.NewQuestionsList(m.conceptId, m.w, func(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
				return m.Update(msg)
			})
		}
	case tea.WindowSizeMsg:
		m.w = msg
		m.help.Width = msg.Width
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func readNotes(conceptId string) string {
	notes, err := ReaderConceptNotes(conceptId)

	if err != nil {
		return fmt.Sprintf("Error reading notes: %s", err)
	}

	return notes
}

func (m conceptModel) View() string {
	helpView := m.help.View(m.keys)
	return fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView(), helpView)
}

func (m conceptModel) headerView() string {
	title := styles.ViewPortTitleStyle.Render(m.title)
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Height(title)))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m conceptModel) footerView() string {
	info := styles.ViewPortInfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func NewConcept(id, title string, w tea.WindowSizeMsg, backHandler BackHandler) (tea.Model, tea.Cmd) {
	// read notes
	content := readNotes(id)
	headerHeight := lipgloss.Height(styles.ViewPortTitleStyle.Render(title))
	footerHeight := lipgloss.Height(styles.ViewPortInfoStyle.Render(""))

	verticalMarginHeight := headerHeight + footerHeight
	availableHeight := w.Height - verticalMarginHeight - helpHeight

	vp := viewport.New(w.Width, availableHeight)
	vp.YPosition = headerHeight
	vp.HighPerformanceRendering = userHighPerformanceRender
	vp.SetContent(content)

	m := conceptModel{
		conceptId: id,
		title:     title,
		help:      help.New(),
		keys:      viewportKeys,
		w:         w,
		viewport:  vp,
		back:      backHandler,
	}

	cmd := tea.SetWindowTitle(title)
	return m, cmd

}
