# Copilot Instructions for headless_unity_build

## Project Overview

This is a Go TUI (Terminal User Interface) application built with Bubble Tea that provides an interactive interface for configuring and executing Unity Android builds in headless mode. The app features two main views: an input configuration view and a build log viewer.

## Build & Run

```bash
# Build the application
go build -o test_app .

# Run the application
./test_app

# Run directly without building
go run .
```

Debug logs are written to `debug.log` in the project root.

## Architecture

### Bubble Tea MVC Pattern

The application follows the Elm Architecture pattern (Model-View-Update) as implemented by Bubble Tea:

- **AppModel** (`model.go`): Root model containing two sub-models and tracks which view is focused
  - **InputModel**: Configuration form for Unity build parameters
  - **BuildLogModel**: Live log output from build command execution

- **View Switching**: Press `Tab` to toggle focus between InputView and LogView. Only the focused view receives keyboard input (except global keys like `q`/`ctrl+c` and `B`).

- **Command Execution**: Uses Go channels for async streaming of stdout/stderr from build commands. The `ExecuteBuildCommand` method spawns a goroutine that pipes output to a channel, which is consumed via `WaitForBuildResponses` to convert channel messages into Bubble Tea Msgs.

### File Responsibilities

- `main.go`: Program entry point, root Update/View logic, global key handling
- `model.go`: Type definitions for all models and message types
- `input_view.go`: Form UI for Unity build configuration with vim-style navigation (j/k for up/down, i to edit)
- `build_log_view.go`: Viewport-based log display with live streaming from exec.Command pipes
- `keys.go`: Key binding definitions (currently defined but not actively used in UI)
- `styles.go`: Lipgloss style constants for focused/selected/blurred states

## Key Conventions

### State Management Pattern

- **focusedIndex vs selectedInput**: `focusedIndex` (0-5) tracks which input has focus (highlighted), `selectedInput` is set only when actively editing (-1 when not editing). This allows vim-style navigation separate from editing mode.

- **Input Updates**: All inputs update on every msg, but only the actively selected input responds to text input. Lists (buildMethod, aabFlag) only update when their selectedInput index matches.

### Message Flow for Async Operations

1. `ExecuteBuildCommand` returns a tea.Cmd that immediately sends a BuildMsg and spawns a goroutine
2. The goroutine streams output to the `sub` channel
3. `WaitForBuildResponses` blocks on the channel and converts messages to BuildMsg or BuildCompleteMsg
4. Each BuildMsg triggers another `WaitForBuildResponses` call to continue listening
5. Special "__DONE__" sentinel value terminates the stream

### Unity Command Construction

The `PrintBuildSettings()` method in `input_view.go` shows the intended Unity command structure:
- Uses Unity headless flags: `-batchmode -nographics -quit`
- Executes a C# method via `-executeMethod AndroidBuildMethod.{BuildUnlock|BuildLock|Build}`
- Passes build parameters after `--` separator
- Currently `ConvertToBuildCmd()` returns empty array (implementation incomplete)

### Style System

Three style states for inputs:
- **noStyle**: Default/blurred (gray)
- **focusedStyle**: When focusedIndex points to it (pink #205)
- **selectedStyle**: When actively editing (teal #90, bold)

Two style states for view containers:
- **modelStyle**: Hidden border
- **focusedModelStyle**: Visible border (blue #69)

## Known Issues / Incomplete Features

- `ConvertToBuildCmd()` currently returns empty slice - build command construction is not implemented
- Key bindings in `keys.go` are defined but not integrated into the UI help system
- Hardcoded Unity path in `PrintBuildSettings()`: `/Unity/Editor/6000.0.58f2/Editor/Unity`
- Build execution is triggered by `B` key (shift+b) but calls `ls` as placeholder instead of actual Unity command
