package practice

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const baseUrl = "https://raw.githubusercontent.com/Thwani47/termilearn-sourcefiles/master"

var (
	spinnerTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	spinnerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
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

type fileDownloadedMsg struct{}
type runTestsMsg struct{ message string }

func downloadFile(concept string) tea.Cmd {
	return func() tea.Msg {
		dir := fmt.Sprintf("practice/concepts/%s", concept)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return errorMsg{err}
		}

		out, err := os.Create(fmt.Sprintf("practice/concepts/%s/main.go", concept))
		if err != nil {
			return errorMsg{err}
		}
		defer out.Close()
		resp, err := http.Get(fmt.Sprintf("%s/%s/practice-questions/%s/main.go", baseUrl, "Go", concept))
		if err != nil {
			return errorMsg{err}
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return errorMsg{fmt.Errorf("Error downloading file: %s", resp.Status)}
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return errorMsg{err}
		}
		return fileDownloadedMsg{}
	}
}

func (p practiceModel) Init() tea.Cmd {
	// TODO: check if file for the concept exists and download it if it does not
	// TODO: this should be return a tea.Cmd (maybe batch the commands)
	// return tea.Batch(p.spinner.Tick, downloadFileFunc)
	// TODO: add help menu for the user (edit (e) , t (test), reset (r), b (back) , q (quit))

	return p.spinner.Tick
}

func (p practiceModel) View() string {

	if p.err != nil {
		return fmt.Sprintf("\n\nError occured: %s", p.err.Error())
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
	return func() tea.Msg {
		cmd := exec.Command("go", "test", fmt.Sprintf("practice/concepts/%s/main.go", concept))

		output, err := cmd.CombinedOutput()
		/*	var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr

			err := cmd.Run()*/
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

	cmd := downloadFile(concept)

	_ = tea.SetWindowTitle("Practice") // specify concept being practices
	return p, cmd
}
