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

type mainModel struct {
	Tabs          []string
	conceptsModel conceptsModel
	activeTab     int
	windowHeight  int
	windowWidth   int
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tea.ClearScreen,
		func() tea.Msg {
			return tea.WindowSizeMsg{
				Width:  80,
				Height: 50,
			}
		})
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if IsQuitting(msg) {
			cmd = tea.Quit
		}

		switch keypress := msg.String(); keypress {
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}

		// delegate key events to teh conceptsModel if the first tab is actie
		if m.activeTab == 0 {
			var model tea.Model
			model, cmd = m.conceptsModel.Update(msg)
			m.conceptsModel = model.(conceptsModel)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	}

	return m, cmd
}

func (m mainModel) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab

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

	var tabContent string

	switch m.activeTab {
	case 0:
		tabContent = m.conceptsModel.View()
	default:
		tabContent = "Coming soon..."
	}
	contentHeight := m.windowHeight - lipgloss.Height(row) - tabDocStyle.GetVerticalPadding()
	doc.WriteString(windowStyle.Height(contentHeight).Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(tabContent))
	return tabDocStyle.Render(doc.String())

}

func NewMainModel() mainModel {
	conceptsModel := NewConceptsModel()
	tabs := []string{"Go Concepts", "Practice Questions", "Interview Questions", "Configuration"}
	return mainModel{
		Tabs:          tabs,
		conceptsModel: conceptsModel,
	}
}
