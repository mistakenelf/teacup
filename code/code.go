package code

import (
	"path/filepath"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/dirfs"
	"github.com/knipferrc/teacup/formatter"
)

type syntaxMsg string
type errorMsg error

const (
	Padding = 1
)

// readFileContentCmd reads the content of the file.
func readFileContentCmd(fileName string) tea.Cmd {
	return func() tea.Msg {
		content, err := dirfs.ReadFileContent(fileName)
		if err != nil {
			return errorMsg(err)
		}

		highlightedContent, err := formatter.Highlight(content, filepath.Ext(fileName), "dracula")
		if err != nil {
			return errorMsg(err)
		}

		return syntaxMsg(highlightedContent)
	}
}

// Bubble represents the properties of a code bubble.
type Bubble struct {
	Viewport           viewport.Model
	Borderless         bool
	Filename           string
	HighlightedContent string
}

// New creates a new instance of code.
func New(borderless bool) Bubble {
	viewPort := viewport.New(0, 0)
	border := lipgloss.NormalBorder()

	if borderless {
		border = lipgloss.HiddenBorder()
	}

	viewPort.Style = lipgloss.NewStyle().
		PaddingLeft(Padding).
		PaddingRight(Padding).
		Border(border)

	return Bubble{
		Viewport:   viewPort,
		Borderless: borderless,
	}
}

func (b Bubble) Init() tea.Cmd {
	return nil
}

// SetFileName sets the content of text.
func (b *Bubble) SetFileName(filename string) {
	b.Filename = filename
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
		b.HighlightedContent = lipgloss.NewStyle().
			Width(b.Viewport.Width).
			Height(b.Viewport.Height).
			Render("Error: " + msg.Error())

		b.Viewport.SetContent(b.HighlightedContent)

		return b, nil
	case tea.WindowSizeMsg:
		b.SetSize(msg.Width, msg.Height)
	}

	if b.Filename != "" {
		cmds = append(cmds, readFileContentCmd(b.Filename))
	}

	b.Viewport, cmd = b.Viewport.Update(msg)
	cmds = append(cmds, cmd)

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
		Border(border)

	return b.Viewport.View()
}
