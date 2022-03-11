package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/code"
)

// Bubble represents the properties of the UI.
type Bubble struct {
	code code.Bubble
}

// New creates a new instance of the UI.
func New() Bubble {
	codeModel := code.New(false, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})

	return Bubble{
		code: codeModel,
	}
}

// Init intializes the UI.
func (b Bubble) Init() tea.Cmd {
	cmd := b.code.SetFileName("code/code.go")

	return cmd
}

// Update handles all UI interactions.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.code.SetSize(msg.Width, msg.Height)

		return b, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	b.code, cmd = b.code.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

// View returns a string representation of the UI.
func (b Bubble) View() string {
	return b.code.View()
}

func main() {
	b := New()
	p := tea.NewProgram(b, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
