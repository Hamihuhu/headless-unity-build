package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func initialModel() AppModel {
	return AppModel{
		focusedView: InputView,
		inputModel:  newInputModel(),
		logModel:    newBuildLogModel(),
	}
}

func (a AppModel) Init() tea.Cmd {
	// return tea.Batch(a.logModel.ExecuteBuildCommand("unityhub", a.inputModel.ConvertToBuildCmd()), WaitForBuildResponses(a.logModel.sub))
	return nil
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// handle global key messages here (e.g., quit)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		// Calculate half width for each view, accounting for borders and padding
		a.viewWidth = (msg.Width / 2) - 4
		a.viewHeight = msg.Height - 4

		// Update viewport size
		a.logModel.viewport.Width = a.viewWidth
		a.logModel.viewport.Height = a.viewHeight

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "tab":
			a.focusedView = (a.focusedView + 1) % 2
		case "B":
			cmds = append(cmds, a.logModel.ExecuteBuildCommand("/Unity/Editor/6000.0.58f2/Editor/Unity", a.inputModel.ConvertToBuildCmd()), WaitForBuildResponses(a.logModel.sub))
		}
	}

	var cmd tea.Cmd
	if a.focusedView == InputView {
		a.inputModel, cmd = a.inputModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	a.logModel, cmd = a.logModel.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a AppModel) View() string {
	inputView := a.inputModel.View()
	logView := a.logModel.View()

	// Apply dynamic sizing
	focusedStyle := focusedModelStyle.Width(a.viewWidth).Height(a.viewHeight)
	unfocusedStyle := modelStyle.Width(a.viewWidth).Height(a.viewHeight)

	if a.focusedView == InputView {
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			focusedStyle.Render(inputView),
			unfocusedStyle.Render(logView),
		)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		unfocusedStyle.Render(inputView),
		focusedStyle.Render(logView),
	)
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func AppendPrintLn(cmd *[]tea.Cmd, log string) {
	*cmd = append(*cmd, tea.Println(log))
}