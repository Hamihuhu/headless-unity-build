package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// type appModel struct {
// 	focusedView int
// 	inputView   InputModel
// 	buildView   BuildSettings
// }

type appModel struct {
	inputModel InputModel
}

func initialModel() appModel {
	return appModel{
		inputModel: newInputModel(),
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
		}
	}

	cmd = a.inputModel.Update(msg)
	return a, cmd
}

func (a appModel) View() string {
	return a.inputModel.View()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}