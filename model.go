package main

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
)

type appModel struct {
	focusedView int
	inputView   inputModel.Model
	buildView   buildSettings.Model
}

type buildSettings struct {
	method      textinput.Model
	buildDir    textinput.Model
	fileName    textinput.Model
	version     textinput.Model
	versionCode textinput.Model
	keystore    textinput.Model
	keyalias    textinput.Model
	aab         bool
}

type inputModel struct {
	focusTextInput int
	textInputs     []textinput.Model
	cursorModel    cursor.Mode
}
