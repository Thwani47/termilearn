package concept

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

const baseUrl = "https://raw.githubusercontent.com/Thwani47/termilearn-sourcefiles/master"

type practiceModel struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	back         BackHandler
	w            tea.WindowSizeMsg
	err          error
}

type clearErrorMsg struct{}
type doneEditingMsg struct{ err error }

func (p practiceModel) Init() tea.Cmd {
	return nil
}

func (p practiceModel) View() string {
	if p.quitting {
		return ""
	}
	return "Practice"
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
	case clearErrorMsg:
		p.err = nil
	}

	return p, nil
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

func NewPractice(w tea.WindowSizeMsg, backhandler BackHandler) (tea.Model, tea.Cmd) {

	p := practiceModel{
		back: backhandler,
		w:    w,
	}

	cmd := tea.SetWindowTitle("Practice") // specify concept being practices
	return p, cmd
}
