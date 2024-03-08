package filetree

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mistakenelf/teacup/filesystem"
)

func (m Model) Init() tea.Cmd {
	return getDirectoryListingCmd(filesystem.CurrentDirectory, true)
}
