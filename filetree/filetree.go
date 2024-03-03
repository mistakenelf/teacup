package filetree

import (
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
	active       bool
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func New(currentDirectory string, active bool) Model {
	fp := filepicker.New()
	fp.CurrentDirectory = currentDirectory

	return Model{
		filepicker:   fp,
		selectedFile: "",
		quitting:     false,
		err:          nil,
		active:       active,
	}
}

func (m Model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case clearErrorMsg:
		m.err = nil
	case tea.WindowSizeMsg:
		m.filepicker.Height = msg.Height
	}

	var cmd tea.Cmd

	if m.active {
		m.filepicker, cmd = m.filepicker.Update(msg)
	}

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selectedFile = path
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""

		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m Model) View() string {
	var s strings.Builder

	if m.quitting {
		return ""
	}

	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	}

	s.WriteString("\n" + m.filepicker.View())

	return s.String()
}
