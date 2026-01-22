package main

import (
	//"fmt"
	// "os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	selectedStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

type InputModel struct {
	buildMethod    list.Model
	buildDirectory textinput.Model
	outputFileName textinput.Model
	versionNumber  textinput.Model
	versionCode    textinput.Model
	aabFlag        list.Model

	focusedIndex  int
	selectedInput int
	cursorModel   cursor.Mode
}

type buildMethodItem string

func (b buildMethodItem) Title() string { return string(b) }
func (b buildMethodItem) Description() string {
	descMap := map[string]string{
		"AndroidBuildMethod.BuildUnlock": "Use BuildUnlock to build the project.",
		"AndroidBuildMethod.BuildLock":   "Use BuildLock to build the project.",
		"AndroidBuildMethod.Build":       "Just build, keep existing settings.",
	}
	return descMap[string(b)]
}
func (b buildMethodItem) FilterValue() string { return string(b) }

func newInputModel() InputModel {
	// Build Method list
	items := []list.Item{
		buildMethodItem("AndroidBuildMethod.BuildUnlock"),
		buildMethodItem("AndroidBuildMethod.BuildLock"),
		buildMethodItem("AndroidBuildMethod.Build"),
	}
	buildMethodList := list.New(items, list.NewDefaultDelegate(), 50, 6)
	buildMethodList.Title = "Select Build Method"
	buildMethodList.SetShowHelp(false)

	// buildDirectory text input
	buildDirInput := textinput.New()
	buildDirInput.Placeholder = "/path/to/build/directory"
	buildDirInput.Prompt = "Build Directory: >"
	buildDirInput.CharLimit = 256
	buildDirInput.Width = 40

	// outputFileName text input
	outputFileNameInput := textinput.New()
	outputFileNameInput.Placeholder = "output.apk"
	outputFileNameInput.Prompt = "Output File Name: >"
	outputFileNameInput.CharLimit = 100
	outputFileNameInput.Width = 40

	// versionNumber text input
	versionNumberInput := textinput.New()
	versionNumberInput.Placeholder = "1.0.0"
	versionNumberInput.Prompt = "Version Number: >"
	versionNumberInput.CharLimit = 20
	versionNumberInput.Width = 20

	// versionCode text input
	versionCodeInput := textinput.New()
	versionCodeInput.Placeholder = "1"
	versionCodeInput.Prompt = "Version Code: >"
	versionCodeInput.CharLimit = 10
	versionCodeInput.Width = 10

	// AAB Flag list
	aabItems := []list.Item{
		buildMethodItem("True"),
		buildMethodItem("False"),
	}
	aabList := list.New(aabItems, list.NewDefaultDelegate(), 50, 2)
	aabList.Title = "Build as AAB?"
	aabList.SetShowHelp(false)

	return InputModel{
		buildMethod:    buildMethodList,
		buildDirectory: buildDirInput,
		outputFileName: outputFileNameInput,
		versionNumber:  versionNumberInput,
		versionCode:    versionCodeInput,
		aabFlag:        aabList,
		focusedIndex:   0,
		selectedInput:  -1,
		cursorModel:    cursor.CursorBlink,
	}
}

func (i *InputModel) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	if i.selectedInput == -1 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "k":
				i.focusNextInput(-1)
			case "j":
				i.focusNextInput(1)
			case "enter":
				i.selectInput()
			}
		}
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				i.cancelInput()
			}
		}
	}

	// update all text inputs
	cmds = append(cmds, i.updateInputs(msg)...)
	return tea.Batch(cmds...)
}

func (i InputModel) View() string {
	var b strings.Builder
	b.WriteString("Android Build Configuration\n\n")
	b.WriteString(i.buildMethod.View() + "\n\n")
	b.WriteString(i.buildDirectory.View() + "\n\n")
	b.WriteString(i.outputFileName.View() + "\n\n")
	b.WriteString(i.versionNumber.View() + "\n\n")
	b.WriteString(i.versionCode.View() + "\n\n")
	b.WriteString(i.aabFlag.View() + "\n\n")

	b.WriteString(helpStyle.Render(strconv.Itoa(i.focusedIndex) + ": Focused Index\n"))
	b.WriteString(helpStyle.Render(strconv.Itoa(i.selectedInput) + ": Selected Input\n"))

	return b.String()
}

func (i *InputModel) focusNextInput(index int) {
	if i.focusedIndex+index < 5 && i.focusedIndex+index >= 0 {
		i.focusedIndex += index
	} else if i.focusedIndex+index < 0 {
		i.focusedIndex = 5
	} else {
		i.focusedIndex = 0
	}

	i.clearStyles()

	// update focused style
	switch i.focusedIndex {
	case 1:
		i.buildDirectory.PromptStyle = focusedStyle
		i.buildDirectory.TextStyle = focusedStyle
	case 2:
		i.outputFileName.PromptStyle = focusedStyle
		i.outputFileName.TextStyle = focusedStyle
	case 3:
		i.versionNumber.PromptStyle = focusedStyle
		i.versionNumber.TextStyle = focusedStyle
	case 4:
		i.versionCode.PromptStyle = focusedStyle
		i.versionCode.TextStyle = focusedStyle
	}
}

func (i *InputModel) selectInput() {
	i.selectedInput = i.focusedIndex

	// focus the selected input
	switch i.selectedInput {
	case 1:
		i.buildDirectory.Focus()
		i.buildDirectory.PromptStyle = focusedStyle
		i.buildDirectory.TextStyle = focusedStyle
	case 2:
		i.outputFileName.Focus()
		i.outputFileName.PromptStyle = focusedStyle
		i.outputFileName.TextStyle = focusedStyle
	case 3:
		i.versionNumber.Focus()
		i.versionNumber.PromptStyle = focusedStyle
		i.versionNumber.TextStyle = focusedStyle
	case 4:
		i.versionCode.Focus()
		i.versionCode.PromptStyle = focusedStyle
		i.versionCode.TextStyle = focusedStyle
	}
}

func (i *InputModel) cancelInput() {
	i.clearStyles()

	switch i.selectedInput {
	case 1:
		i.buildDirectory.Blur()
		i.buildDirectory.PromptStyle = focusedStyle
		i.buildDirectory.TextStyle = focusedStyle
	case 2:
		i.outputFileName.Blur()
		i.outputFileName.PromptStyle = focusedStyle
		i.outputFileName.TextStyle = focusedStyle
	case 3:
		i.versionNumber.Blur()
		i.versionNumber.PromptStyle = focusedStyle
		i.versionNumber.TextStyle = focusedStyle
	case 4:
		i.versionCode.Blur()
		i.versionCode.PromptStyle = focusedStyle
		i.versionCode.TextStyle = focusedStyle
	}

	i.selectedInput = -1
}

func (i *InputModel) updateInputs(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	if i.selectedInput == 0 {
		i.buildMethod, cmd = i.buildMethod.Update(msg)
		cmds = append(cmds, cmd)
	}

	if i.selectedInput == 5 {
		i.aabFlag, cmd = i.aabFlag.Update(msg)
		cmds = append(cmds, cmd)
	}

	i.buildDirectory, cmd = i.buildDirectory.Update(msg)
	cmds = append(cmds, cmd)

	i.outputFileName, cmd = i.outputFileName.Update(msg)
	cmds = append(cmds, cmd)

	i.versionNumber, cmd = i.versionNumber.Update(msg)
	cmds = append(cmds, cmd)

	i.versionCode, cmd = i.versionCode.Update(msg)
	cmds = append(cmds, cmd)

	return cmds
}

func (i *InputModel) clearStyles() {
	i.buildDirectory.PromptStyle = noStyle
	i.buildDirectory.TextStyle = noStyle

	i.outputFileName.PromptStyle = noStyle
	i.outputFileName.TextStyle = noStyle

	i.versionNumber.PromptStyle = noStyle
	i.versionNumber.TextStyle = noStyle

	i.versionCode.PromptStyle = noStyle
	i.versionCode.TextStyle = noStyle
}
