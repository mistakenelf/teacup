// Package statusbar provides an statusbar bubble which can render
// four different status sections
package statusbar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// Height represents the height of the statusbar.
const Height = 1

// ColorConfig
type ColorConfig struct {
	Foreground lipgloss.AdaptiveColor
	Background lipgloss.AdaptiveColor
}

// Bubble represents the properties of the statusbar.
type Bubble struct {
	Width              int
	Height             int
	FirstColumn        string
	SecondColumn       string
	ThirdColumn        string
	FourthColumn       string
	FirstColumnColors  ColorConfig
	SecondColumnColors ColorConfig
	ThirdColumnColors  ColorConfig
	FourthColumnColors ColorConfig
}

// New creates a new instance of the statusbar.
func New(firstColumnColors, secondColumnColors, thirdColumnColors, fourthColumnColors ColorConfig) Bubble {
	return Bubble{
		FirstColumnColors:  firstColumnColors,
		SecondColumnColors: secondColumnColors,
		ThirdColumnColors:  thirdColumnColors,
		FourthColumnColors: fourthColumnColors,
	}
}

// SetSize sets the width of the statusbar.
func (b *Bubble) SetSize(width int) {
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

// SetColors sets the colors of the 4 columns.
func (b *Bubble) SetColors(firstColumnColors, secondColumnColors, thirdColumnColors, fourthColumnColors ColorConfig) {
	b.FirstColumnColors = firstColumnColors
	b.SecondColumnColors = secondColumnColors
	b.ThirdColumnColors = thirdColumnColors
	b.FourthColumnColors = fourthColumnColors
}

// View returns a string representation of a statusbar.
func (b Bubble) View() string {
	width := lipgloss.Width

	firstColumn := lipgloss.NewStyle().
		Foreground(b.FirstColumnColors.Foreground).
		Background(b.FirstColumnColors.Background).
		Padding(0, 1).
		Height(Height).
		Render(truncate.StringWithTail(b.FirstColumn, 30, "..."))

	thirdColumn := lipgloss.NewStyle().
		Foreground(b.ThirdColumnColors.Foreground).
		Background(b.ThirdColumnColors.Background).
		Align(lipgloss.Right).
		Padding(0, 1).
		Height(Height).
		Render(b.ThirdColumn)

	fourthColumn := lipgloss.NewStyle().
		Foreground(b.FourthColumnColors.Foreground).
		Background(b.FourthColumnColors.Background).
		Padding(0, 1).
		Height(Height).
		Render(b.FourthColumn)

	secondColumn := lipgloss.NewStyle().
		Foreground(b.SecondColumnColors.Foreground).
		Background(b.SecondColumnColors.Background).
		Padding(0, 1).
		Height(Height).
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
