package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

// initial state of the model
func initialModel() model {
	return model{
		choices:  []string{"apple", "banana", "cherry", "date", "elderberry"},
		selected: make(map[int]struct{}),
	}
}

// do nothing here to init
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

}