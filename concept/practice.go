package concept

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

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
	return p.filepicker.Init()
}

func (p practiceModel) View() string {
	if p.quitting {
		return ""
	}

	var sb strings.Builder

	sb.WriteString("\n ")
	if p.err != nil {
		sb.WriteString(p.filepicker.Styles.DisabledFile.Render(p.err.Error()))
	} else if p.selectedFile == "" {
		sb.WriteString("Open file") // TODO specify file name maybe? what about concepts that need multiple files? Will we have that
	} else {
		sb.WriteString("selected file  " + p.filepicker.Styles.Selected.Render(p.selectedFile))
	}

	sb.WriteString("\n\n" + p.filepicker.View() + "\n")
	return sb.String()
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

	var cmd tea.Cmd
	p.filepicker, cmd = p.filepicker.Update(msg)

	if didSelect, path := p.filepicker.DidSelectDisabledFile(msg); didSelect {
		p.err = errors.New(path + "is not valid")
		p.selectedFile = ""
		return p, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return p, cmd
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

func clearErrorAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func NewPractice(w tea.WindowSizeMsg, backhandler BackHandler) (tea.Model, tea.Cmd) {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".go"}
	fp.CurrentDirectory, _ = os.Getwd()

	p := practiceModel{
		filepicker: fp,
		back:       backhandler,
		w:          w,
	}

	cmd := tea.SetWindowTitle("Practice") // specify concept being practices
	return p, cmd
}
