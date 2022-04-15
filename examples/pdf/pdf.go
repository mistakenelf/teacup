package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/pdf"
)

// Bubble represents the properties of the UI.
type Bubble struct {
	pdf pdf.Bubble
}

// New creates a new instance of the UI.
func New() Bubble {
	pdfModel := pdf.New(true, true, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})

	return Bubble{
		pdf: pdfModel,
	}
}

// Init intializes the UI.
func (b Bubble) Init() tea.Cmd {
	return nil
}

// Update handles all UI interactions.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.pdf.SetSize(msg.Width, msg.Height)
		cmds = append(cmds, b.pdf.SetFileName("./examples/pdf/sample-pdf-file.pdf"))

		return b, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	b.pdf, cmd = b.pdf.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

// View returns a string representation of the UI.
func (b Bubble) View() string {
	return b.pdf.View()
}

func main() {
	b := New()
	p := tea.NewProgram(b, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
