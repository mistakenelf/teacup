// Package pdf provides an pdf bubble which can render
// pdf files as strings.
package pdf

import (
	"bytes"
	"errors"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ledongthuc/pdf"
)

type renderPDFMsg string
type errorMsg error

const (
	padding = 1
)

// Bubble represents the properties of a pdf bubble.
type Bubble struct {
	Viewport    viewport.Model
	BorderColor lipgloss.AdaptiveColor
	Active      bool
	Borderless  bool
	FileName    string
}

// ReadPdf reads a PDF file given a name.
func ReadPdf(name string) (string, error) {
	file, reader, err := pdf.Open(name)
	if err != nil {
		return "", errors.Unwrap(err)
	}

	defer func() {
		if e := file.Close(); e != nil {
			err = e
		}
	}()

	buf := new(bytes.Buffer)
	buffer, err := reader.GetPlainText()

	if err != nil {
		return "", errors.Unwrap(err)
	}

	_, err = buf.ReadFrom(buffer)
	if err != nil {
		return "", errors.Unwrap(err)
	}

	return buf.String(), nil
}

// renderPDFCmd reads the content of a PDF and returns its content as a string.
func renderPDFCmd(filename string) tea.Cmd {
	return func() tea.Msg {
		pdfContent, err := ReadPdf(filename)
		if err != nil {
			return errorMsg(err)
		}

		return renderPDFMsg(pdfContent)
	}
}

// New creates a new instance of a PDF.
func New(active, borderless bool, borderColor lipgloss.AdaptiveColor) Bubble {
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

	return Bubble{
		Viewport:    viewPort,
		Borderless:  borderless,
		BorderColor: borderColor,
	}
}

// SetBorderless sets weather or not to show the border.
func (b *Bubble) SetBorderless(borderless bool) {
	b.Borderless = borderless
}

// Init initializes the PDF bubble.
func (b Bubble) Init() tea.Cmd {
	return nil
}

// SetFileName sets current file to render, this
// returns a cmd which will render the pdf.
func (b *Bubble) SetFileName(filename string) tea.Cmd {
	b.FileName = filename

	return renderPDFCmd(filename)
}

// SetBorderColor sets the current color of the border.
func (b *Bubble) SetBorderColor(color lipgloss.AdaptiveColor) {
	b.BorderColor = color
}

// SetSize sets the size of the bubble.
func (b *Bubble) SetSize(w, h int) {
	b.Viewport.Width = w - b.Viewport.Style.GetHorizontalFrameSize()
	b.Viewport.Height = h - b.Viewport.Style.GetVerticalFrameSize()

	border := lipgloss.NormalBorder()

	if b.Borderless {
		border = lipgloss.HiddenBorder()
	}

	b.Viewport.Style = lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border).
		BorderForeground(b.BorderColor)
}

// SetIsActive sets if the bubble is currently active.
func (b *Bubble) SetIsActive(active bool) {
	b.Active = active
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
	case renderPDFMsg:
		pdfContent := lipgloss.NewStyle().
			Width(b.Viewport.Width).
			Height(b.Viewport.Height).
			Render(string(msg))

		b.Viewport.SetContent(pdfContent)

		return b, nil
	case errorMsg:
		b.FileName = ""
		b.Viewport.SetContent(msg.Error())

		return b, nil
	}

	if b.Active {
		b.Viewport, cmd = b.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return b, tea.Batch(cmds...)
}

// View returns a string representation of the markdown bubble.
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
