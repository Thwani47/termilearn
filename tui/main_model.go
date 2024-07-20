package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appState int

const (
	viewConceptList appState = iota
	viewConcept
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
	Tabs         []string
	TabContent   []string
	activeTab    int
	conceptsList tea.Model
	concept      conceptModel
	state        appState
	width        int
	height       int
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case viewConceptList:
		m.conceptsList, cmd = m.conceptsList.Update(msg)
	case viewConcept:
		var model tea.Model
		model, cmd = m.concept.Update(msg)
		m.concept = model.(conceptModel)
	}
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
	case conceptSelectedMessage:
		m.state = viewConcept
		m.concept = NewConceptModel(msg.id, msg.choice)
		var model tea.Model
		model, _ = m.concept.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		m.concept = model.(conceptModel)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, cmd
}

func (m mainModel) View() string {
	/*switch m.state {
	case viewConceptList:
		return m.conceptsList.View()
	case viewConcept:
		return m.concept.View()
	default:
		return "I don't think we should be here"
	}*/

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
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.activeTab]))
	return tabDocStyle.Render(doc.String())

}

func NewMainModel() mainModel {
	conceptsList := NewConceptListModel()
	tabs := []string{"Go Concepts", "Practice Questions", "Interview Questions", "Configuration"}
	tabContent := []string{"Work in Progress ", "Coming Soon", "Coming Soon", "Coming Soon"}
	return mainModel{
		Tabs:         tabs,
		TabContent:   tabContent,
		conceptsList: conceptsList,
		state:        viewConceptList,
	}
}
