package help

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Constants used for the help bubble.
const (
	Padding  = 1
	KeyWidth = 12
)

// Entry represents a single entry in the help bubble.
type Entry struct {
	Key         string
	Description string
}

// Bubble represents the struct of the help bubble.
type Bubble struct {
	Title       string
	Entries     []Entry
	Viewport    viewport.Model
	BorderColor lipgloss.AdaptiveColor
	Borderless  bool
}

// generateHelpScreen generates the help text based on the title and entries.
func (b Bubble) generateHelpScreen(entries []Entry, title string) string {
	helpScreen := ""

	for _, content := range entries {
		keyText := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"}).
			Width(KeyWidth).
			Render(content.Key)

		descriptionText := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"}).
			Render(content.Description)

		row := lipgloss.JoinHorizontal(lipgloss.Top, keyText, descriptionText)
		helpScreen += fmt.Sprintf("%s\n", row)
	}

	welcomeText := lipgloss.NewStyle().Bold(true).
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Border(lipgloss.NormalBorder()).
		Padding(0, 1).
		Italic(true).
		BorderBottom(true).
		BorderTop(false).
		BorderRight(false).
		BorderLeft(false).
		Render(title)

	return lipgloss.NewStyle().
		Width(b.Viewport.Width).
		Height(b.Viewport.Height).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			welcomeText,
			helpScreen,
		))
}

// New creates a new help bubble.
func New(
	borderColor lipgloss.AdaptiveColor,
	title string,
	entries []Entry,
	borderless bool,
) Bubble {
	viewPort := viewport.New(0, 0)
	border := lipgloss.NormalBorder()

	if borderless {
		border = lipgloss.HiddenBorder()
	}

	viewPort.Style = lipgloss.NewStyle().
		PaddingLeft(Padding).
		PaddingRight(Padding).
		Border(border).
		BorderForeground(borderColor)

	return Bubble{
		Viewport:   viewPort,
		Entries:    entries,
		Title:      title,
		Borderless: borderless,
	}
}

// SetSize sets the size of the help bubble.
func (b *Bubble) SetSize(w, h int) {
	b.Viewport.Width = w - b.Viewport.Style.GetHorizontalFrameSize()
	b.Viewport.Height = h - b.Viewport.Style.GetVerticalFrameSize()

	b.Viewport.SetContent(b.generateHelpScreen(b.Entries, b.Title))
}

// SetBorderColor sets the border color of the help bubble.
func (b *Bubble) SetBorderColor(color lipgloss.AdaptiveColor) {
	b.BorderColor = color
}

// Update handles updating the help bubble.
func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.SetSize(msg.Width, msg.Height)
	}

	b.Viewport, cmd = b.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

// View returns a string representation of the help bubble.
func (b Bubble) View() string {
	border := lipgloss.NormalBorder()

	if b.Borderless {
		border = lipgloss.HiddenBorder()
	}

	b.Viewport.Style = lipgloss.NewStyle().
		PaddingLeft(Padding).
		PaddingRight(Padding).
		Border(border).
		BorderForeground(b.BorderColor)

	return b.Viewport.View()
}
