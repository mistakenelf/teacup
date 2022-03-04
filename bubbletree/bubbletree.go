package bubbletree

import "github.com/charmbracelet/lipgloss"

type Bubble struct{}

func (b Bubble) View() string {
	return lipgloss.NewStyle().Bold(true).Render("Bubbletree")
}
