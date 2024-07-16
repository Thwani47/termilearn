package tui

import (
	"fmt"
	"strings"

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

type conceptModel struct {
	Title    string
	ready    bool
	viewport viewport.Model
}

func (m conceptModel) Init() tea.Cmd {
	return nil
}

func (m conceptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.viewport.SetContent(m.Title)
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

func (m conceptModel) View() string {
	if !m.ready {
		return "\n Initializing... " + m.Title
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m conceptModel) headerView() string {
	title := viewPortTitleStyle.Render("Mr. Pager")
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Height(title)))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m conceptModel) footerView() string {
	info := viewPortInfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
func NewConceptModel(title string) conceptModel {
	model := conceptModel{
		Title: title,
	}

	return model
}
