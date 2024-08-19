package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/Thwani47/termilearn/common/keys"
	"github.com/Thwani47/termilearn/common/styles"
	"github.com/Thwani47/termilearn/concept"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	concepts = []list.Item{
		listItem{id: "hello-world", title: "Hello World", description: "Learn how to print the classical 'Hello World' message in Go"},
		listItem{id: "values", title: "Values", description: "Learn about the basic types in Go"},
		listItem{id: "variables", title: "Variables", description: "Learn how to declare and initialize variables in Go"},
		listItem{id: "constants", title: "Constants", description: "Learn how to declare and initialize constants in Go"},
		listItem{id: "for-loop", title: "For", description: "Learn how to loop in Go"},
		listItem{id: "if-else", title: "If/Else", description: "Learn how to use conditional statements in Go"},
		listItem{id: "switch", title: "Switch", description: "Learn how to use switch statements in Go"},
	}
)

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
	list        list.Model
	keys        keys.ConceptListKeyMap
	backHandler BackHandler
	w           tea.WindowSizeMsg
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

	fn := styles.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.SelectedItemStyle.Render("> " + strings.Join(s, " "))
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
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			return m.backHandler(m.w)
		case key.Matches(msg, m.keys.Choose):
			i, ok := m.list.SelectedItem().(listItem)
			if ok {
				return concept.NewConcept(i.id, i.title, m.w, func(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
					return m.Update(msg)
				})
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.w = msg
		h, v := styles.DocStyle.GetFrameSize()
		availableHeight := msg.Height - v - 2
		m.list.SetSize(msg.Width-h, availableHeight)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m conceptListModel) View() string {
	return styles.DocStyle.Render(m.list.View())
}

func NewConceptList(width int, height int, backHandler BackHandler) (tea.Model, tea.Cmd) {
	keys := keys.ConceptListKeys
	l := list.New(concepts, itemDelegate{}, 0, 0)

	l.Title = "Select Go Concept"
	l.Styles.Title = styles.TitleStyle
	l.Styles.PaginationStyle = styles.PaginationStyle
	l.SetShowStatusBar(false)
	_, v := styles.DocStyle.GetFrameSize()
	availableHeight := height - v - 2
	l.SetSize(width, availableHeight)
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Choose}
	}
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Choose}
	}

	cmd := l.NewStatusMessage("")

	conceptListModel := conceptListModel{
		list: l,
		w: tea.WindowSizeMsg{
			Width:  width,
			Height: height,
		},
		keys:        keys,
		backHandler: backHandler,
	}

	return conceptListModel, cmd
}
