// Package code implements a code bubble which renders syntax highlighted
// source code based on a filename.
package code

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/dirfs"
)

type syntaxMsg string
type errorMsg error

// Constants used throughout.
const (
	Padding = 1
)

// Highlight returns a syntax highlighted string of text.
func Highlight(content, extension, syntaxTheme string) (string, error) {
	buf := new(bytes.Buffer)
	if err := quick.Highlight(buf, content, extension, "terminal256", syntaxTheme); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return buf.String(), nil
}

// readFileContentCmd reads the content of the file.
func readFileContentCmd(fileName, syntaxTheme string) tea.Cmd {
	return func() tea.Msg {
		content, err := dirfs.ReadFileContent(fileName)
		if err != nil {
			return errorMsg(err)
		}

		highlightedContent, err := Highlight(content, filepath.Ext(fileName), syntaxTheme)
		if err != nil {
			return errorMsg(err)
		}

		return syntaxMsg(highlightedContent)
	}
}

// Bubble represents the properties of a code bubble.
type Bubble struct {
	Viewport           viewport.Model
	BorderColor        lipgloss.AdaptiveColor
	Borderless         bool
	Active             bool
	Filename           string
	HighlightedContent string
	SyntaxTheme        string
}

// New creates a new instance of code.
func New(active, borderless bool, borderColor lipgloss.AdaptiveColor) Bubble {
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
		Viewport:    viewPort,
		Borderless:  borderless,
		Active:      active,
		BorderColor: borderColor,
	}
}

// Init initializes the code bubble.
func (b Bubble) Init() tea.Cmd {
	return nil
}

// SetFileName sets current file to highlight.
func (b *Bubble) SetFileName(filename string) tea.Cmd {
	b.Filename = filename

	return readFileContentCmd(filename, b.SyntaxTheme)
}

// SetIsActive sets if the bubble is currently active.
func (b *Bubble) SetIsActive(active bool) {
	b.Active = active
}

// SetBorderColor sets the current color of the border.
func (b *Bubble) SetBorderColor(color lipgloss.AdaptiveColor) {
	b.BorderColor = color
}

// SetSyntaxTheme sets the syntax theme of the rendered code.
func (b *Bubble) SetSyntaxTheme(theme string) {
	b.SyntaxTheme = theme
}

// SetBorderless sets weather or not to show the border.
func (b *Bubble) SetBorderless(borderless bool) {
	b.Borderless = borderless
}

// SetSize sets the size of the bubble.
func (b *Bubble) SetSize(w, h int) {
	b.Viewport.Width = w - b.Viewport.Style.GetHorizontalFrameSize()
	b.Viewport.Height = h - b.Viewport.Style.GetVerticalFrameSize()

	b.Viewport.SetContent(lipgloss.NewStyle().
		Width(b.Viewport.Width).
		Height(b.Viewport.Height).
		Render(b.HighlightedContent))
}

// GotoTop jumps to the top of the viewport.
func (b *Bubble) GotoTop() {
	b.Viewport.GotoTop()
}

// Update handles updating the UI of a code bubble.
func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case syntaxMsg:
		b.Filename = ""
		b.HighlightedContent = lipgloss.NewStyle().
			Width(b.Viewport.Width).
			Height(b.Viewport.Height).
			Render(string(msg))

		b.Viewport.SetContent(b.HighlightedContent)

		return b, nil
	case errorMsg:
		b.Filename = ""
		b.HighlightedContent = lipgloss.NewStyle().
			Width(b.Viewport.Width).
			Height(b.Viewport.Height).
			Render("Error: " + msg.Error())

		b.Viewport.SetContent(b.HighlightedContent)

		return b, nil
	}

	if b.Active {
		b.Viewport, cmd = b.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return b, tea.Batch(cmds...)
}

// View returns a string representation of the code bubble.
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
