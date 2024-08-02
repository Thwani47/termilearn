package practice

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	spinnerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	errorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render
	helpHeight       = 6
)

type BackHandler func(tea.WindowSizeMsg) (tea.Model, tea.Cmd)

type practiceModel struct {
	concept           string
	err               error
	quitting          bool
	help              help.Model
	keys              keyMap
	back              BackHandler
	w                 tea.WindowSizeMsg
	isDownloadingFile bool
	spinner           spinner.Model
	testOutput        string
}

type doneEditingMsg struct{ err error }
type errorMsg struct{ err error }

func (e errorMsg) Error() string {
	return e.err.Error()
}

func (p practiceModel) Init() tea.Cmd {
	return p.spinner.Tick
}

func (p practiceModel) View() string {

	if p.err != nil {
		return fmt.Sprintf("\n\n%s\n\n", errorStyle(p.err.Error()))
	}

	if p.isDownloadingFile {
		return fmt.Sprintf("\n%s %s\n\n", p.spinner.View(), spinnerTextStyle("Setting up..."))
	}

	if p.testOutput != "" {
		return fmt.Sprintf("\n%s\n\n%s\n\n", spinnerTextStyle("Tests ran successfully"), p.testOutput)
	}

	helpView := p.help.View(p.keys)

	return fmt.Sprintf("\n%s\n\n%s\n\n", spinnerTextStyle("Ready for practice..."), helpView)

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
		case key.Matches(msg, p.keys.Open):
			return p, practiceConcept(p.concept)
		case key.Matches(msg, p.keys.Test):
			return p, testConcept(p.concept)
		}
	case tea.WindowSizeMsg:
		p.w = msg
	case errorMsg:
		p.err = msg.err
		p.isDownloadingFile = false
		return p, nil
	case fileDownloadedMsg:
		p.isDownloadingFile = false
		return p, nil
	case runTestsMsg:
		p.testOutput = msg.message
		return p, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		p.spinner, cmd = p.spinner.Update(msg)
		return p, cmd
	}

	return p, p.spinner.Tick
}

func testConcept(concept string) tea.Cmd {
	// TODO: how can we capture the tests and show them as a list to the user (maybe show a checkmark next to the test that passed and a cross next to the test that failed)
	return func() tea.Msg {
		cmd := exec.Command("go", "test", fmt.Sprintf("practice/concepts/%s/%s_test.go", concept, concept))

		output, err := cmd.CombinedOutput()
		if err != nil {
			return errorMsg{fmt.Errorf(fmt.Sprint(err) + ": " + string(output))}
		}

		return runTestsMsg{string(output)}
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

func NewPractice(concept string, w tea.WindowSizeMsg, backhandler BackHandler) (tea.Model, tea.Cmd) {

	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = spinnerStyle
	p := practiceModel{
		concept:           concept,
		back:              backhandler,
		w:                 w,
		help:              help.New(),
		keys:              keys,
		isDownloadingFile: true,
		spinner:           s,
	}

	cmd := getPracticeFiles(concept)

	_ = tea.SetWindowTitle(fmt.Sprintf("Practice: %s", concept))
	return p, cmd
}
