package tui

import (
	"strings"

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
)

type mainModelState int

const (
	viewTabs mainModelState = iota
	viewConcepts
	viewPracticeQuestions
	viewInterviewQuestions
	viewConfiguration
)

type mainModel struct {
	tabs          []string
	tabContent    []string
	conceptsModel conceptsModel
	activeTab     int
	state         mainModelState
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if IsQuitting(msg) {
			cmd = tea.Quit
		}
		if m.state == viewTabs {
			switch keypress := msg.String(); keypress {
			case "right", "n", "tab":
				m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
				return m, nil
			case "left", "p", "shift+tab":
				m.activeTab = max(m.activeTab-1, 0)
				return m, nil
			case "enter":
				switch m.activeTab {
				case 0:
					m.state = viewConcepts
				case 1:
					m.state = viewPracticeQuestions
				case 2:
					m.state = viewInterviewQuestions
				case 3:
					m.state = viewConfiguration
				}
				return m, nil
			}
		} else if m.state == viewConcepts {
			// delegate key events to the conceptsModel
			var subCmd tea.Cmd
			var model tea.Model
			model, subCmd = m.conceptsModel.Update(msg)
			m.conceptsModel = model.(conceptsModel)
			return m, subCmd
		}
	}

	return m, cmd
}

func (m mainModel) View() string {
	if m.state == viewConcepts {
		return m.conceptsModel.View()
	}

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
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.tabContent[m.activeTab]))
	return tabDocStyle.Render(doc.String())

}

func NewMainModel() mainModel {
	conceptsModel := NewConceptsModel()
	tabs := []string{"Go Concepts", "Practice Questions", "Interview Questions", "Configuration"}
	tabContent := []string{"View Concepts", "View Go practice questions", "View Go interview questions", "Configure termilearn"}
	return mainModel{
		tabs:          tabs,
		tabContent:    tabContent,
		conceptsModel: conceptsModel,
		state:         viewTabs,
	}
}
