package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("90")).Bold(true)
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle       = lipgloss.NewStyle()
	helpStyle     = blurredStyle
)

var (
	modelStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2)

	focusedModelStyle = lipgloss.NewStyle().
				Align(lipgloss.Left).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("205")).
				Padding(1, 2)

	listFocusedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205"))

	listSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("90")).
				Bold(true)
)