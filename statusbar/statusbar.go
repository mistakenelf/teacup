// Package statusbar provides an statusbar bubble which can render
// four different status sections
package statusbar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

const statusBarHeight = 1

// Bubble represents the properties of the statusbar.
type Bubble struct {
	Width        int
	Height       int
	FirstColumn  string
	SecondColumn string
	ThirdColumn  string
	FourthColumn string
}

// SetSize sets the width of the statusbar.
func (b *Bubble) SetSize(width) {
	b.Width = width
}

// Update updates the size of the statusbar.
func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.SetSize(msg.Width)
	}

	return b, nil
}


// SetContent sets the content of the statusbar.
func (b *Bubble) SetContent(firstColumn, secondColumn, thirdColumn, fourthColumn string) {
  b.FirstColumn = firstColumn
  b.SecondColumn = secondColumn
  b.ThirdColumn = thirdColumn
  b.FourthColumn = fourthColumn
}

func (b Bubble) View() string {
  width := lipgloss.Width

	firstColumn := lipgloss.NewStyle().
		Padding(0, 1).
		Height(statusBarHeight).
		Render(truncate.StringWithTail(b.FirstColumn, 30, "..."))

	thirdColumn := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Padding(0, 1).
		Height(statusBarHeight).
		Render(b.ThirdColumn)

	fourthColumn := lipgloss.NewStyle().
		Padding(0, 1).
		Height(statusBarHeight).
		Render(b.FourthColumn)

	secondColumn := lipgloss.NewStyle().
		Padding(0, 1).
		Height(statusBarHeight).
		Width(b.Width - width(firstColumn) - width(thirdColumn) - width(fourthColumn)).
		Render(truncate.StringWithTail(
			b.SecondColumn,
			uint(b.Width-width(firstColumn)-width(thirdColumn)-width(fourthColumn)-3),
			"..."),
		)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		firstColumn,
		secondColumn,
		thirdColumn,
		fourthColumn,
	)
}
