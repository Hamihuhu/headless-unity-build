package main

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
)

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

type appModel struct {
	focusTextInput int
	textInputs     []textinput.Model
	cursorModel    cursor.Mode
}