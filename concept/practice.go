package concept

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const baseUrl = "https://raw.githubusercontent.com/Thwani47/termilearn-sourcefiles/master"

var (
	spinnerTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	spinnerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
)

type practiceModel struct {
	err               error
	quitting          bool
	back              BackHandler
	w                 tea.WindowSizeMsg
	isDownloadingFile bool
	spinner           spinner.Model
}

type doneEditingMsg struct{ err error }

func (p practiceModel) Init() tea.Cmd {
	// TODO: check if file for the concept exists and download it if it does not
	// TODO: this should be return a tea.Cmd (maybe batch the commands)
	// return tea.Batch(p.spinner.Tick, downloadFileFunc)
	// TODO: add help menu for the user (edit (e) , t (test), reset (r), b (back) , q (quit))
	return p.spinner.Tick
}

func (p practiceModel) View() string {

	return fmt.Sprintf("\n%s %s\n\n", p.spinner.View(), spinnerTextStyle("Downloading..."))
}

func (p practiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			p.quitting = true
			return p, tea.Quit
		case "b":
			return p.back(p.w)
		}
	case tea.WindowSizeMsg:
		p.w = msg
	case spinner.TickMsg:
		var cmd tea.Cmd
		p.spinner, cmd = p.spinner.Update(msg)
		return p, cmd
	}

	return p, p.spinner.Tick
}

func practiceConcept(concept string) tea.Cmd {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		editor = "vim"
	}
	file := fmt.Sprintf("concept/practice/%s/%s.go", concept, concept)
	command := exec.Command(editor, file)

	return tea.ExecProcess(command, func(err error) tea.Msg {
		return doneEditingMsg{err}
	})
}

// TODO: we need to pass in the title of the concept
func NewPractice(w tea.WindowSizeMsg, backhandler BackHandler) (tea.Model, tea.Cmd) {

	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = spinnerStyle
	p := practiceModel{
		back:    backhandler,
		w:       w,
		spinner: s,
	}

	cmd := tea.SetWindowTitle("Practice") // specify concept being practices
	return p, cmd
}
