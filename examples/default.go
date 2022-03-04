package main

import (
	"log"

	"github.com/knipferrc/bubbletree/bubbletree"

	tea "github.com/charmbracelet/bubbletea"
)

type Bubble struct {
	bubbletree bubbletree.Bubble
}

func New() Bubble {
	return Bubble{}
}

func (b Bubble) Init() tea.Cmd {
	return nil
}

func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return b, tea.Batch(cmds...)
}

func (b Bubble) View() string {
	return b.bubbletree.View()
}

func main() {
	b := New()
	p := tea.NewProgram(b, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
