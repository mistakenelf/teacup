package filetree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/filesystem"
)

const (
	thousand    = 1000
	ten         = 10
	fivePercent = 0.0499
)

var (
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
)

type KeyMap struct {
	Down key.Binding
	Up   key.Binding
}

type DirectoryItem struct {
	name             string
	details          string
	path             string
	extension        string
	isDirectory      bool
	currentDirectory string
}

type getDirectoryListingMsg []DirectoryItem
type errorMsg error

type Model struct {
	viewport viewport.Model
	cursor   int
	files    []DirectoryItem
	active   bool
	keyMap   KeyMap
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Down: key.NewBinding(key.WithKeys("j", "down", "ctrl+n"), key.WithHelp("j", "down")),
		Up:   key.NewBinding(key.WithKeys("k", "up", "ctrl+p"), key.WithHelp("k", "up")),
	}
}

func New() Model {
	viewPort := viewport.New(0, 0)

	return Model{
		viewport: viewPort,
		cursor:   0,
		active:   true,
		keyMap:   DefaultKeyMap(),
	}
}

// ConvertBytesToSizeString converts a byte count to a human readable string.
func ConvertBytesToSizeString(size int64) string {
	if size < thousand {
		return fmt.Sprintf("%dB", size)
	}

	suffix := []string{
		"K", // kilo
		"M", // mega
		"G", // giga
		"T", // tera
		"P", // peta
		"E", // exa
		"Z", // zeta
		"Y", // yotta
	}

	curr := float64(size) / thousand
	for _, s := range suffix {
		if curr < ten {
			return fmt.Sprintf("%.1f%s", curr-fivePercent, s)
		} else if curr < thousand {
			return fmt.Sprintf("%d%s", int(curr), s)
		}
		curr /= thousand
	}

	return ""
}

// SetSize sets the size of the bubble.
func (m *Model) SetSize(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h
}

// SetIsActive sets if the bubble is currently active.
func (m *Model) SetIsActive(active bool) {
	m.active = active
}

// GotoTop jumps to the top of the viewport.
func (m *Model) GotoTop() {
	m.viewport.GotoTop()
}

// getDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func getDirectoryListingCmd(directoryName string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		var directoryItems []DirectoryItem

		if directoryName == filesystem.HomeDirectory {
			directoryName, err = filesystem.GetHomeDirectory()
			if err != nil {
				return errorMsg(err)
			}
		}

		directoryInfo, err := os.Stat(directoryName)
		if err != nil {
			return errorMsg(err)
		}

		if !directoryInfo.IsDir() {
			return nil
		}

		files, err := filesystem.GetDirectoryListing(directoryName, showHidden)
		if err != nil {
			return errorMsg(err)
		}

		err = os.Chdir(directoryName)
		if err != nil {
			return errorMsg(err)
		}

		workingDirectory, err := filesystem.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err)
		}

		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				continue
			}

			status := fmt.Sprintf("%s %s %s",
				fileInfo.ModTime().Format("2006-01-02 15:04:05"),
				fileInfo.Mode().String(),
				ConvertBytesToSizeString(fileInfo.Size()))

			directoryItems = append(directoryItems, DirectoryItem{
				name:             file.Name(),
				details:          status,
				path:             filepath.Join(workingDirectory, file.Name()),
				extension:        filepath.Ext(fileInfo.Name()),
				isDirectory:      fileInfo.IsDir(),
				currentDirectory: workingDirectory,
			})
		}

		return getDirectoryListingMsg(directoryItems)
	}
}

func (m Model) Init() tea.Cmd {
	return getDirectoryListingCmd(filesystem.CurrentDirectory, true)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
	case getDirectoryListingMsg:
		if msg != nil {
			m.files = msg
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Up):
			m.cursor--
		case key.Matches(msg, m.keyMap.Down):
			m.cursor++
		}
	}

	return m, nil
}

func (m Model) View() string {
	var fileList strings.Builder

	for i, file := range m.files {
		if i == m.cursor {
			fileList.WriteString(selectedItemStyle.Render(file.name) + "\n")
		} else {
			fileList.WriteString(file.name + "\n")
		}
	}

	return fileList.String()
}
