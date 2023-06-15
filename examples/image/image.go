package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/image"
)

// model represents the properties of the UI.
type model struct {
	image image.Model
}

// New creates a new instance of the UI.
func New() model {
	imageModel := image.New(true, true, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})

	return model{
		image: imageModel,
	}
}

// Init intializes the UI.
func (b model) Init() tea.Cmd {
	return nil
}

// Update handles all UI interactions.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.image.SetSize(msg.Width, msg.Height)
		cmds = append(cmds, m.image.SetFileName("examples/image/bubbletea.png"))

		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	m.image, cmd = m.image.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the UI.
func (m model) View() string {
	return m.image.View()
}

func main() {
	b := New()
	p := tea.NewProgram(b, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
