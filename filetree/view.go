// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View returns a string representation of a filetree.
func (m Model) View() string {
	var inputView string

	switch m.state {
	case idleState:
		inputView = ""
	case createFileState, createDirectoryState, renameItemState:
		inputView = m.input.View()
	case deleteItemState:
		inputView = "Are you sure you want to delete? (y/n)"
	case moveItemState:
		inputView = fmt.Sprintf("Currently moving %s", m.itemToMove.shortName)
	default:
		inputView = ""
	}

	return bubbleStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			m.list.View(),
			inputStyle.Render(inputView),
		))
}
