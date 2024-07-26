package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	tabDocStyle       = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	helpHeight        = 6
)

type mainModel struct {
	activeTab int
	width     int
	height    int
	keys      tabsKeyMap
	help      help.Model
	tabs      []Tab
}

type Tab struct {
	title   string
	content string
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Next):
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case key.Matches(msg, m.keys.Prev):
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Select):
			switch m.activeTab {
			case 0:
				return NewConceptList(m.width, m.height, func(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
					return m.Update(msg)
				})
			case 1:
				return m, nil
			case 2:
				return m, nil

			case 3:
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
	}
	return m, cmd
}

func (m mainModel) View() string {

	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab

		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}

		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t.title))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.tabs[m.activeTab].content))
	helpView := m.help.View(m.keys)
	doc.WriteString(fmt.Sprintf("\n%s", helpView))
	return tabDocStyle.Render(doc.String())

}

func NewMainModel() mainModel {
	tabs := []Tab{
		{title: "Go Concepts", content: "View Concepts"},
		{title: "Practice Questions", content: "View Go practice questions"},
		{title: "Interview Questions", content: "View Go interview questions"},
		{title: "Configuration", content: "Configure termilearn"},
	}
	return mainModel{
		tabs: tabs,
		help: help.New(),
		keys: tabKeys,
	}
}
