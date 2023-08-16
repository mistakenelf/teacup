package filetree

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mistakenelf/teacup/dirfs"
)

// Init initializes the filetree with files from the current directory.
func (m Model) Init() tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if m.startDir == "" {
		cmd = getDirectoryListingCmd(dirfs.CurrentDirectory, m.showHidden, m.showIcons)
	} else {
		cmd = getDirectoryListingCmd(m.startDir, m.showHidden, m.showIcons)
	}

	cmds = append(cmds, cmd, textinput.Blink)

	return tea.Batch(cmds...)
}
