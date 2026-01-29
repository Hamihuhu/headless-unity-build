package main

import (
	"bufio"
	// "fmt"
	// "io"
	// "os"
	"os/exec"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func newBuildLogModel() BuildLogModel {
	vp := viewport.New(100, 20)
	vp.SetContent("Build log will appear here...\n")
	return BuildLogModel{
		sub:       make(chan string),
		logString: "",
		viewport:  vp,
	}
}

func (b BuildLogModel) Update(msg tea.Msg) (BuildLogModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case BuildMsg:
		b.logString += string(msg) + "\n"
		b.viewport.SetContent(b.logString)
		b.viewport.GotoBottom()
		cmds = append(cmds, WaitForBuildResponses(b.sub))

	case BuildErrorMsg:
		b.logString += "Error: " + msg.err.Error() + "\n"

	case BuildCompleteMsg:
		b.logString += "Done.\n"
	}
	var cmd tea.Cmd
	b.viewport, cmd = b.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

func (b BuildLogModel) View() string {
	return b.viewport.View()
}

func (b *BuildLogModel) ExecuteBuildCommand(executablePath string, buildCmd []string) tea.Cmd {
	return func() tea.Msg {
		go func() {
			b.sub <- "Starting build command: " + executablePath
			cmd := exec.Command(executablePath, buildCmd...)

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				b.sub <- "ERROR: " + err.Error()
				b.sub <- "__DONE__"
				return
			}

			stderr, err := cmd.StderrPipe()
			if err != nil {
				b.sub <- "ERROR: " + err.Error()
				b.sub <- "__DONE__"
				return
			}

			if err := cmd.Start(); err != nil {
				b.sub <- "ERROR: " + err.Error()
				b.sub <- "__DONE__"
				return
			}

			// Read both stdout and stderr
			done := make(chan bool)
			
			go func() {
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					b.sub <- "[OUT] " + scanner.Text()
				}
				done <- true
			}()

			go func() {
				scanner := bufio.NewScanner(stderr)
				for scanner.Scan() {
					b.sub <- "[ERR] " + scanner.Text()
				}
				done <- true
			}()

			<-done
			<-done

			if err := cmd.Wait(); err != nil {
				b.sub <- "Command error: " + err.Error()
			}

			b.sub <- "__DONE__"
		}()

		return BuildMsg("Command queued: " + executablePath)
	}
}

func (b *BuildLogModel) appendLog(logLine string) {
	b.logString += logLine + "\n"
	b.viewport.SetContent(b.logString)
	b.viewport.GotoBottom()
}

func WaitForBuildResponses(sub chan string) tea.Cmd {
	return func() tea.Msg {
		msg := <-sub
		if msg == "__DONE__" {
			return BuildCompleteMsg{}
		}
		return BuildMsg(msg)
	}
}
