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

// Model represents the properties of a pdf bubble.
type Model struct {
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
func New(active, borderless bool, borderColor lipgloss.AdaptiveColor) Model {
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

	return Model{
		Viewport:    viewPort,
		Borderless:  borderless,
		BorderColor: borderColor,
	}
}

// SetBorderless sets weather or not to show the border.
func (m *Model) SetBorderless(borderless bool) {
	m.Borderless = borderless
}

// Init initializes the PDF bubble.
func (m Model) Init() tea.Cmd {
	return nil
}

// SetFileName sets current file to render, this
// returns a cmd which will render the pdf.
func (m *Model) SetFileName(filename string) tea.Cmd {
	m.FileName = filename

	return renderPDFCmd(filename)
}

// SetBorderColor sets the current color of the border.
func (m *Model) SetBorderColor(color lipgloss.AdaptiveColor) {
	m.BorderColor = color
}

// SetSize sets the size of the bubble.
func (m *Model) SetSize(w, h int) {
	m.Viewport.Width = w
	m.Viewport.Height = h

	border := lipgloss.NormalBorder()

	if m.Borderless {
		border = lipgloss.HiddenBorder()
	}

	m.Viewport.Style = lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border).
		BorderForeground(m.BorderColor)
}

// SetIsActive sets if the bubble is currently active.
func (m *Model) SetIsActive(active bool) {
	m.Active = active
}

// GotoTop jumps to the top of the viewport.
func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

// Update handles updating the UI of a code bubble.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case renderPDFMsg:
		pdfContent := lipgloss.NewStyle().
			Width(m.Viewport.Width).
			Height(m.Viewport.Height).
			Render(string(msg))

		m.Viewport.SetContent(pdfContent)

		return m, nil
	case errorMsg:
		m.FileName = ""
		m.Viewport.SetContent(msg.Error())

		return m, nil
	}

	if m.Active {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the markdown bubble.
func (m Model) View() string {
	border := lipgloss.NormalBorder()

	if m.Borderless {
		border = lipgloss.HiddenBorder()
	}

	m.Viewport.Style = lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border).
		BorderForeground(m.BorderColor)

	return m.Viewport.View()
}
