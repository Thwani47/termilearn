package practice

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	questionListStyle = lipgloss.NewStyle().Margin(1, 2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	spinnerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	spinnerTextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render

	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type QuestionListItem struct {
	title       string
	description string
}

type questionList struct {
	questionsList          []QuestionWrapper
	questions              list.Model
	isDownloadingQuestions bool
	err                    error
	concept                string
	spinner                spinner.Model
	w                      tea.WindowSizeMsg
	back                   BackHandler
}

type questionDelegate struct{}

func (q QuestionListItem) Title() string                            { return q.title }
func (q QuestionListItem) Description() string                      { return q.description }
func (q QuestionListItem) FilterValue() string                      { return q.title }
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

func (q questionList) Init() tea.Cmd { return nil }

func (q questionList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		q.w = msg
	case errorMsg:
		q.isDownloadingQuestions = false
		return q, nil
	case fileDownloadedMsg:
		q.isDownloadingQuestions = false
		q.questionsList = msg.questions
		q.err = msg.err

	case tea.KeyMsg:
		if msg.String() == "q" {
			return q, tea.Quit
		}
	}

	return q, q.spinner.Tick
}

func (q questionList) View() string {
	if q.isDownloadingQuestions {
		return fmt.Sprintf("\n%s %s\n\n", q.spinner.View(), spinnerTextStyle("Setting up..."))
	}

	if q.err != nil {
		return errorStyle(fmt.Sprintf("\n\n%s\n\n", q.err.Error()))
	}

	if q.questionsList != nil {
		q.createListView(q.questionsList)
		return questionListStyle.Render(q.questions.View())
	}

	return errorStyle(fmt.Sprintf("\n\nNo questions found for concept: %s\n\n", q.concept))
}

func (q *questionList) createListView(questions []QuestionWrapper) {
	items := make([]list.Item, len(questions))
	for i, q := range questions {
		var title, description string
		switch q.QuestionType {
		case "mcq":
			title = q.MCQQuestion.Title
			description = q.MCQQuestion.QuestionText
		case "edit":
			title = q.EditQuestion.Title
			description = fmt.Sprintf("File: %s, TestFile: %s", q.EditQuestion.File, q.EditQuestion.TestFile)
		}
		items[i] = QuestionListItem{title: title, description: description}
	}

	l := list.New(items, questionDelegate{}, 0, 0)
	l.SetShowStatusBar(false)
	l.Title = q.concept
	_, v := questionListStyle.GetFrameSize()
	l.SetSize(q.w.Width, q.w.Height-v-2)
	q.questions = l
}

func NewQuestionsList(concept string, w tea.WindowSizeMsg, backhandler BackHandler) (tea.Model, tea.Cmd) {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = spinnerStyle
	m := questionList{
		concept:                concept,
		w:                      w,
		back:                   backhandler,
		isDownloadingQuestions: true,
		spinner:                s,
	}

	cmd := getPracticeFiles(concept)
	return m, cmd
}
