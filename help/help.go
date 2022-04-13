// Package help implements a help bubble which can be used
// to display help information such as keymaps.
package help

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 1
	keyWidth = 12
)

type TitleColor struct {
	Background lipgloss.AdaptiveColor
	Foreground lipgloss.AdaptiveColor
}

// Entry represents a single entry in the help bubble.
type Entry struct {
	Key         string
	Description string
}

// Bubble represents the properties of a help bubble.
type Bubble struct {
	Viewport    viewport.Model
	Entries     []Entry
	BorderColor lipgloss.AdaptiveColor
	Title       string
	TitleColor  TitleColor
	Active      bool
	Borderless  bool
}

// generateHelpScreen generates the help text based on the title and entries.
func generateHelpScreen(title string, titleColor TitleColor, entries []Entry, width, height int) string {
	helpScreen := ""

	for _, content := range entries {
		keyText := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"}).
			Width(keyWidth).
			Render(content.Key)

		descriptionText := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"}).
			Render(content.Description)

		row := lipgloss.JoinHorizontal(lipgloss.Top, keyText, descriptionText)
		helpScreen += fmt.Sprintf("%s\n", row)
	}

	titleText := lipgloss.NewStyle().Bold(true).
		Background(titleColor.Background).
		Foreground(titleColor.Foreground).
		Border(lipgloss.NormalBorder()).
		Padding(0, 1).
		Italic(true).
		BorderBottom(true).
		BorderTop(false).
		BorderRight(false).
		BorderLeft(false).
		Render(title)

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			titleText,
			helpScreen,
		))
}

// New creates a new instance of a help bubble.
func New(
	active, borderless bool,
	title string,
	titleColor TitleColor,
	borderColor lipgloss.AdaptiveColor,
	entries []Entry,
) Bubble {
	viewPort := viewport.New(0, 0)
	border := lipgloss.NormalBorder()

	if borderless {
		border = lipgloss.HiddenBorder()
	}

	viewPort.Style = lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border).
		BorderForeground(borderColor)

	viewPort.SetContent(generateHelpScreen(title, titleColor, entries, 0, 0))

	return Bubble{
		Viewport:    viewPort,
		Entries:     entries,
		Title:       title,
		Active:      active,
		Borderless:  borderless,
		BorderColor: borderColor,
		TitleColor:  titleColor,
	}
}

// SetSize sets the size of the help bubble.
func (b *Bubble) SetSize(w, h int) {
	b.Viewport.Width = w - b.Viewport.Style.GetHorizontalFrameSize()
	b.Viewport.Height = h - b.Viewport.Style.GetVerticalFrameSize()

	b.Viewport.SetContent(generateHelpScreen(b.Title, b.TitleColor, b.Entries, b.Viewport.Width, b.Viewport.Height))
}

// SetBorderColor sets the current color of the border.
func (b *Bubble) SetBorderColor(color lipgloss.AdaptiveColor) {
	b.BorderColor = color
}

// SetIsActive sets if the bubble is currently active.
func (b *Bubble) SetIsActive(active bool) {
	b.Active = active
}

// GotoTop jumps to the top of the viewport.
func (b *Bubble) GotoTop() {
	b.Viewport.GotoTop()
}

// SetTitleColor sets the color of the title.
func (b *Bubble) SetTitleColor(color TitleColor) {
	b.TitleColor = color

	b.Viewport.SetContent(generateHelpScreen(b.Title, b.TitleColor, b.Entries, b.Viewport.Width, b.Viewport.Height))
}

// SetBorderless sets weather or not to show the border.
func (b *Bubble) SetBorderless(borderless bool) {
	b.Borderless = borderless
}

// Update handles UI interactions with the help bubble.
func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if b.Active {
		b.Viewport, cmd = b.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return b, tea.Batch(cmds...)
}

// View returns a string representation of the help bubble.
func (b Bubble) View() string {
	border := lipgloss.NormalBorder()

	if b.Borderless {
		border = lipgloss.HiddenBorder()
	}

	b.Viewport.Style = lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border).
		BorderForeground(b.BorderColor)

	return b.Viewport.View()
}
