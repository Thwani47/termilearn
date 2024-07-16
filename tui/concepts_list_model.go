package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle          = lipgloss.NewStyle().Margin(1, 2)
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("#FFFDF5")).Background(lipgloss.Color("#25A065")).Padding(0, 1)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	concepts          = []list.Item{
		listItem{id: "hello-world", title: "Hello World", description: "Learn how to print the classical 'Hello World' message in Go"},
		listItem{id: "values", title: "Values", description: "Learn about the basic types in Go"},
		listItem{id: "variables", title: "Variables", description: "Learn how to declare and initialize variables in Go"},
		listItem{id: "constants", title: "Constants", description: "Learn how to declare and initialize constants in Go"},
		listItem{id: "for-loop", title: "For", description: "Learn how to loop in Go"},
		listItem{id: "if-else", title: "If/Else", description: "Learn how to use conditional statements in Go"},
		listItem{id: "switch", title: "Switch", description: "Learn how to use switch statements in Go"},
	}
)

type conceptSelectedMessage struct {
	id     string
	choice string
}

type listItem struct {
	id          string
	title       string
	description string
}

func (li listItem) Title() string {
	return li.title
}

func (li listItem) Description() string {
	return li.description
}

func (li listItem) FilterValue() string {
	return li.title
}

type conceptListModel struct {
	list list.Model
}

func (c conceptListModel) Init() tea.Cmd {
	return nil
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(listItem)

	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (m conceptListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		if IsQuitting(msg) {
			return m, tea.Quit
		}
		switch msg.String() {
		case "enter":
			i, ok := m.list.SelectedItem().(listItem)
			if ok {
				return m, func() tea.Msg {
					return conceptSelectedMessage{id: i.id, choice: i.title}
				}
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m conceptListModel) View() string {
	return docStyle.Render(m.list.View())
}

func NewConceptListModel() conceptListModel {
	conceptsListModel := conceptListModel{
		list: list.New(concepts, itemDelegate{}, 0, 0),
	}

	conceptsListModel.list.Title = "Go Concepts"
	conceptsListModel.list.Styles.Title = titleStyle
	conceptsListModel.list.Styles.PaginationStyle = paginationStyle
	conceptsListModel.list.SetShowStatusBar(false)

	return conceptsListModel
}
