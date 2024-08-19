package practice

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Thwani47/termilearn/common/keys"
	"github.com/Thwani47/termilearn/common/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type BackHandler func(tea.WindowSizeMsg) (tea.Model, tea.Cmd)

type practiceModel struct {
	question QuestionWrapper
	err      error
	quitting bool
	help     help.Model
	keys     keys.PracticeKeyMap
	back     BackHandler
	w        tea.WindowSizeMsg
	tests    []testResult
}

type doneEditingMsg struct{ err error }
type errorMsg struct{ err error }

func (e errorMsg) Error() string {
	return e.err.Error()
}

type runTestsMsg struct {
	err   error
	tests []testResult
}

func (r runTestsMsg) Error() string {
	return r.err.Error()
}

type testResult struct {
	name         string
	passed       bool
	errorMessage string
}

func (p practiceModel) Init() tea.Cmd {
	return nil
}

func (p practiceModel) View() string {
	if len(p.tests) > 0 {
		var resultView string

		for _, test := range p.tests {
			if test.passed {
				resultView += fmt.Sprintf("✅ %s\n", test.name)
			} else {
				resultView += fmt.Sprintf("❌ %s\n%s\n", test.name, test.errorMessage)
			}

		}
		return fmt.Sprintf("\n%s\n\n%s\n\n", styles.SpinnerTextStyle("Tests results: "), resultView)

	}

	helpView := p.help.View(p.keys)

	return fmt.Sprintf("\n%s\n\n\n\n%s\n\n", styles.SpinnerTextStyle("practice..."), helpView)

}

func (p practiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, p.keys.Quit):
			p.quitting = true
			return p, tea.Quit
		case key.Matches(msg, p.keys.Back):
			return p.back(p.w)
		case key.Matches(msg, p.keys.Help):
			p.help.ShowAll = !p.help.ShowAll
			// case key.Matches(msg, p.keys.Open):
			// 	return p, practiceConcept(p.concept)
			// case key.Matches(msg, p.keys.Test):
			// 	return p, testConcept(p.concept)
		}
	case tea.WindowSizeMsg:
		p.w = msg
	case errorMsg:
		p.err = msg.err
		return p, nil
	case runTestsMsg:
		if msg.err != nil {
			p.err = msg.err
		}
		p.tests = msg.tests
		return p, nil

	}

	return p, nil
}

func testConcept(concept string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("go", "test", "-json", fmt.Sprintf("practice/concepts/%s/%s_test.go", concept, concept))

		output, err := cmd.CombinedOutput()
		if err != nil {
			tests := parseTestOutput(string(output))
			return runTestsMsg{tests: tests, err: fmt.Errorf(string(output))}
		}

		tests := parseTestOutput(string(output))

		return runTestsMsg{tests: tests}
	}
}

func practiceConcept(concept string) tea.Cmd {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		editor = "vim"
	}
	file := fmt.Sprintf("practice/concepts/%s/main.go", concept)
	command := exec.Command(editor, file)

	return tea.ExecProcess(command, func(err error) tea.Msg {
		return doneEditingMsg{err}
	})
}

func NewPractice(question QuestionWrapper, w tea.WindowSizeMsg, backhandler BackHandler) (tea.Model, tea.Cmd) {
	p := practiceModel{
		question: question,
		back:     backhandler,
		w:        w,
		help:     help.New(),
		keys:     keys.PracticeKeys,
	}

	var title string

	if question.QuestionType == "mcq" {
		title = question.MCQQuestion.Title
	} else {
		title = question.EditQuestion.Title
	}

	cmd := tea.SetWindowTitle(fmt.Sprintf("Practice: %s", title))
	return p, cmd
}
