package tui

import (
	"fmt"
	"strings"

	"github.com/Thwani47/termilearn/helpers"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const userHighPerformanceRender = false

var (
	viewPortTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	viewPortInfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type singleConceptModel struct {
	conceptId string
	title     string
	ready     bool
	viewport  viewport.Model
}

func (m singleConceptModel) Init() tea.Cmd {
	return nil
}

func (m singleConceptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if IsQuitting(msg) {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = userHighPerformanceRender
			m.viewport.SetContent(readNotes(m.conceptId))
			m.ready = true
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if userHighPerformanceRender {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}

	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func readNotes(conceptId string) string {
	notes, err := helpers.ReaderConceptNotes(conceptId)

	if err != nil {
		return fmt.Sprintf("Error reading notes: %s", err)
	}

	return notes
}

func (m singleConceptModel) View() string {
	if !m.ready {
		return "\n Initializing... " + m.title
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m singleConceptModel) headerView() string {
	title := viewPortTitleStyle.Render(m.title)
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Height(title)))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m singleConceptModel) footerView() string {
	info := viewPortInfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func NewSingleConceptModel(id string, title string) singleConceptModel {
	model := singleConceptModel{
		conceptId: id,
		title:     title,
	}

	return model
}
