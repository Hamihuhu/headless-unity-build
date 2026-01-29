package main

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

// //////////////////////////////////////////////
type FocusedView int

const (
	InputView FocusedView = iota
	LogView
)

type AppModel struct {
	focusedView  FocusedView
	inputModel   InputModel
	logModel     BuildLogModel
	width        int
	height       int
	viewWidth    int
	viewHeight   int
}

// //////////////////////////////////////////////
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

// //////////////////////////////////////////////
type BuildLogModel struct {
	logString string
	sub       chan string
	viewport  viewport.Model
}

type BuildMsg string
type BuildCompleteMsg struct{}
type BuildErrorMsg struct{ err error }
