// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/knipferrc/teacup/dirfs"
)

// Init initializes the filetree with files from the current directory.
func (b Bubble) Init() tea.Cmd {
	return getDirectoryListingCmd(dirfs.CurrentDirectory, true)
}
