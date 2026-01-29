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
)

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
	// Create custom delegate with matching colors
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = noStyle
	delegate.Styles.SelectedTitle = listFocusedStyle
	delegate.Styles.NormalDesc = blurredStyle
	delegate.Styles.SelectedDesc = listFocusedStyle

	// Build Method list
	items := []list.Item{
		buildMethodItem("AndroidBuildMethod.BuildUnlock"),
		buildMethodItem("AndroidBuildMethod.BuildLock"),
		buildMethodItem("AndroidBuildMethod.Build"),
	}
	buildMethodList := list.New(items, delegate, 50, 6)
	buildMethodList.Title = "Select Build Method"
	buildMethodList.SetShowHelp(false)

	// buildDirectory text input
	buildDirInput := textinput.New()
	buildDirInput.Placeholder = "/path/to/build/directory"
	buildDirInput.Prompt = "Build Directory: "
	buildDirInput.CharLimit = 256
	buildDirInput.Width = 40

	// outputFileName text input
	outputFileNameInput := textinput.New()
	outputFileNameInput.Placeholder = "output.apk"
	outputFileNameInput.Prompt = "Output File Name: "
	outputFileNameInput.CharLimit = 100
	outputFileNameInput.Width = 40

	// versionNumber text input
	versionNumberInput := textinput.New()
	versionNumberInput.Placeholder = "1.0.0"
	versionNumberInput.Prompt = "Version Number: "
	versionNumberInput.CharLimit = 20
	versionNumberInput.Width = 20

	// versionCode text input
	versionCodeInput := textinput.New()
	versionCodeInput.Placeholder = "1"
	versionCodeInput.Prompt = "Version Code: "
	versionCodeInput.CharLimit = 10
	versionCodeInput.Width = 10

	// AAB Flag list (reuse delegate)
	aabItems := []list.Item{
		buildMethodItem("True"),
		buildMethodItem("False"),
	}
	aabList := list.New(aabItems, delegate, 50, 2)
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

func (i InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	var cmds []tea.Cmd

	// Don't consume shift+b - let the global handler in main.go handle it
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "shift+b" {
		return i, nil
	}

	if i.selectedInput == -1 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "k":
				i.focusNextInput(-1)
			case "j":
				i.focusNextInput(1)
			case "i":
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
	return i, tea.Batch(cmds...)
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
	b.WriteString(i.PrintBuildSettings() + "\n")

	return b.String()
}

func (i *InputModel) focusNextInput(index int) {
	if i.focusedIndex+index <= 5 && i.focusedIndex+index >= 0 {
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
		i.buildDirectory.PromptStyle = selectedStyle
		i.buildDirectory.TextStyle = selectedStyle
	case 2:
		i.outputFileName.Focus()
		i.outputFileName.PromptStyle = selectedStyle
		i.outputFileName.TextStyle = selectedStyle
	case 3:
		i.versionNumber.Focus()
		i.versionNumber.PromptStyle = selectedStyle
		i.versionNumber.TextStyle = selectedStyle
	case 4:
		i.versionCode.Focus()
		i.versionCode.PromptStyle = selectedStyle
		i.versionCode.TextStyle = selectedStyle
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

func (i InputModel) PrintBuildSettings() string {
	var b strings.Builder
	b.WriteString("/Unity/Editor/6000.0.58f2/Editor/Unity")
	b.WriteString(" -batchmode")
	b.WriteString(" -nographics")
	b.WriteString(" -quit")
	b.WriteString(" -projectPath " + i.buildDirectory.Value())
	b.WriteString(" -executeMethod " + i.buildMethod.SelectedItem().FilterValue())
	b.WriteString(" -logFile -")
	b.WriteString(" --")
	b.WriteString(" buildPath=" + i.buildDirectory.Value())
	b.WriteString(" fileName=" + i.outputFileName.Value())
	b.WriteString(" version=" + i.versionNumber.Value())
	b.WriteString(" versionCode=" + i.versionCode.Value())
	b.WriteString(" aabCheck=" + strings.ToLower(i.aabFlag.SelectedItem().FilterValue()))
	return b.String()
}

// func (i InputModel) ConvertToBuildCmd() []string {
// 	var cmd []string
// 	cmd = append(cmd, "-batchmode")
// 	cmd = append(cmd, "-nographics")
// 	cmd = append(cmd, "-quit")
// 	cmd = append(cmd, "-projectPath", i.buildDirectory.Value())
// 	cmd = append(cmd, "-executeMethod", i.buildMethod.SelectedItem().FilterValue())
// 	cmd = append(cmd, "-logFile", "-")
// 	cmd = append(cmd, "--")
// 	cmd = append(cmd, "buildPath="+i.buildDirectory.Value())
// 	cmd = append(cmd, "fileName="+i.outputFileName.Value())
// 	cmd = append(cmd, "version="+i.versionNumber.Value())
// 	cmd = append(cmd, "versionCode="+i.versionCode.Value())
// 	cmd = append(cmd, "aabCheck="+strings.ToLower(i.aabFlag.SelectedItem().FilterValue()))
// 	return cmd
// }

func (i InputModel) ConvertToBuildCmd() []string {
	var cmd []string
	cmd = append(cmd, "-batchmode")
	cmd = append(cmd, "-nographics")
	cmd = append(cmd, "-quit")
	cmd = append(cmd, "-projectPath", "/Unity/Projects/NW3_Build/source")
	cmd = append(cmd, "-executeMethod", "AndroidBuildMethod.BuildUnlock")
	cmd = append(cmd, "-logFile", "-")
	cmd = append(cmd, "--")
	cmd = append(cmd, "buildPath=/Unity/Projects/NW3_Build/build")
	cmd = append(cmd, "fileName=test.apk")
	cmd = append(cmd, "version=60")
	cmd = append(cmd, "versionCode=1.60.2")
	cmd = append(cmd, "aabCheck=false")
	return cmd
}