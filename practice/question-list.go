package practice

import (
	"fmt"
	"io"
	"strings"

	"github.com/Thwani47/termilearn/common/keys"
	"github.com/Thwani47/termilearn/common/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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
	keys                   keys.QuestionListKeyMap
	help                   help.Model
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

	fn := styles.ItemStyle.Render

	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (q questionList) Init() tea.Cmd {
	return tea.Batch(q.spinner.Tick, getPracticeFiles(q.concept))
}

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
		q.createListView(q.questionsList)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, q.keys.Quit):
			return q, tea.Quit
		case key.Matches(msg, q.keys.Back):
			return q.back(q.w)
		case key.Matches(msg, q.keys.Help):
			q.help.ShowAll = !q.help.ShowAll
		}
	}

	var cmd tea.Cmd
	q.questions, cmd = q.questions.Update(msg)
	return q, cmd
}

func (q questionList) View() string {
	if q.isDownloadingQuestions {
		return fmt.Sprintf("\n%s %s\n\n", q.spinner.View(), styles.SpinnerTextStyle("Setting up..."))
	}

	if q.err != nil {
		return styles.ErrorStyle(fmt.Sprintf("\n\n%s\n\n", q.err.Error()))
	}

	if q.questionsList != nil {
		return styles.DocStyle.Render(q.questions.View())
	}

	return styles.ErrorStyle(fmt.Sprintf("\n\nNo questions found for concept: %s\n\n", q.concept))
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

	q.questions.SetItems(items)
}

func NewQuestionsList(concept string, w tea.WindowSizeMsg, backhandler BackHandler) (tea.Model, tea.Cmd) {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = styles.SpinnerStyle

	l := list.New([]list.Item{}, questionDelegate{}, w.Width, w.Height)
	l.SetShowStatusBar(false)
	_, v := styles.DocStyle.GetFrameSize()
	l.SetSize(w.Width, w.Height-v-2)
	l.Title = concept

	m := questionList{
		concept:                concept,
		w:                      w,
		back:                   backhandler,
		isDownloadingQuestions: true,
		help:                   help.New(),
		keys:                   keys.QuestionListKeys,
		spinner:                s,
		questions:              l,
	}

	return m, tea.Batch(m.spinner.Tick, getPracticeFiles(concept))
}
