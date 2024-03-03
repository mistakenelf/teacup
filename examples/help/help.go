package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/help"
)

// model represents the properties of the UI.
type model struct {
	help help.Model
}

// New create a new instance of the UI.
func New() model {
	helpModel := help.New(
		true,
		"Help",
		help.TitleColor{
			Background: lipgloss.AdaptiveColor{Light: "62", Dark: "62"},
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffffs"},
		},
		[]help.Entry{
			{Key: "ctrl+c", Description: "Exit FM"},
			{Key: "j/up", Description: "Move up"},
			{Key: "k/down", Description: "Move down"},
			{Key: "h/left", Description: "Go back a directory"},
			{Key: "l/right", Description: "Read file or enter directory"},
			{Key: "p", Description: "Preview directory"},
			{Key: "G", Description: "Jump to bottom"},
			{Key: "~", Description: "Go to home directory"},
			{Key: ".", Description: "Toggle hidden files"},
			{Key: "y", Description: "Copy file path to clipboard"},
			{Key: "Z", Description: "Zip currently selected tree item"},
			{Key: "U", Description: "Unzip currently selected tree item"},
			{Key: "n", Description: "Create new file"},
			{Key: "N", Description: "Create new directory"},
			{Key: "ctrl+d", Description: "Delete currently selected tree item"},
			{Key: "M", Description: "Move currently selected tree item"},
			{Key: "enter", Description: "Process command"},
			{Key: "E", Description: "Edit currently selected tree item"},
			{Key: "C", Description: "Copy currently selected tree item"},
			{Key: "esc", Description: "Reset FM to initial state"},
			{Key: "tab", Description: "Toggle between boxes"},
		},
	)

	return model{
		help: helpModel,
	}
}

// Init initializes the application.
func (m model) Init() tea.Cmd {
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
		m.help.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI.
func (m model) View() string {
	return m.help.View()
}

func main() {
	b := New()
	p := tea.NewProgram(b, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
