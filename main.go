package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())

	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
)

type appModel struct {
	focusedView int
	inputView   InputModel
	logView     textarea.Model
}

func initialModel() appModel {
	return appModel{
		focusedView: 0,
		inputView:   newInputModel(),
		logView:     textarea.New(),
	}
}

func (a appModel) Init() tea.Cmd {
	return nil
}

func (a appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// handle global key messages here (e.g., quit)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "tab":
			a.focusedView = (a.focusedView + 1) % 2
		}
	}

	switch a.focusedView {
	case 0:
		cmd = a.inputView.Update(msg)
	case 1:
		a.logView, cmd = a.logView.Update(msg)
	}

	cmd = a.inputView.Update(msg)
	return a, cmd
}

func (a appModel) View() string {
	var s string
	if a.focusedView == 0 {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(fmt.Sprintf("%4s", a.inputView.View())), modelStyle.Render(a.logView.View()))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%4s", a.inputView.View())), focusedModelStyle.Render(a.logView.View()))
	}

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}