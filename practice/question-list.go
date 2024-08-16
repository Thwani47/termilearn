package practice

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	questionListStyle = lipgloss.NewStyle().Margin(1, 2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type QuestionListItem struct {
	title       string
	description string
}

func (q QuestionListItem) Title() string       { return q.title }
func (q QuestionListItem) Description() string { return q.description }
func (q QuestionListItem) FilterValue() string { return q.title }

type questionDelegate struct{}

func (qd questionDelegate) Height() int                             { return 1 }
func (qd questionDelegate) Spacing() int                            { return 0 }
func (qd questionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (qd questionDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(QuestionListItem)

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

type questionList struct {
	questions list.Model
}

func (q questionList) Init() tea.Cmd { return nil }

func (q questionList) View() string { return questionListStyle.Render(q.questions.View()) }

func (q questionList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := questionListStyle.GetFrameSize()
		q.questions.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if msg.String() == "q" {
			return q, tea.Quit
		}
	}

	var cmd tea.Cmd
	q.questions, cmd = q.questions.Update(msg)
	return q, cmd
}

func NewQuestionsList(questions []QuestionWrapper, w tea.WindowSizeMsg, backhandler BackHandler) tea.Model {
	items := make([]list.Item, len(questions))
	for i, q := range questions {
		var title, description string
		switch q.QuestionType {
		case "mcq":
			title = q.MCQQuestion.QuestionType
			description = q.MCQQuestion.QuestionText
		case "edit":
			title = q.EditQuestion.QuestionType
			description = fmt.Sprintf("File: %s, TestFile: %s", q.EditQuestion.File, q.EditQuestion.TestFile)
		}
		items[i] = QuestionListItem{title: title, description: description}

	}

	l := list.New(items, questionDelegate{}, 0, 0)
	l.SetShowStatusBar(false)
	_, v := questionListStyle.GetFrameSize()
	l.SetSize(w.Width, w.Height-v-2)

	m := questionList{
		questions: l,
	}
	return m
}
