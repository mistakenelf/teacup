package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mistakenelf/teacup/code"
)

// model represents the properties of the UI.
type model struct {
	code code.Model
}

// New creates a new instance of the UI.
func New() model {
	codeModel := code.New(true)

	return model{
		code: codeModel,
	}
}

// Init intializes the UI.
func (m model) Init() tea.Cmd {
	cmd := m.code.SetFileName("code/code.go")

	return cmd
}

// Update handles all UI interactions.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.code.SetSize(msg.Width, msg.Height)

		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	m.code, cmd = m.code.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the UI.
func (m model) View() string {
	return m.code.View()
}

func main() {
	b := New()
	p := tea.NewProgram(b, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
