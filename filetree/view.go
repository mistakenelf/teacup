// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View returns a string representation of a filetree.
func (b Bubble) View() string {
	var inputView string

	switch b.state {
	case idleState:
		inputView = ""
	case createFileState, createDirectoryState, renameItemState:
		inputView = b.input.View()
	case deleteItemState:
		inputView = "Are you sure you want to delete? (y/n)"
	case moveItemState:
		inputView = fmt.Sprintf("Currently moving %s", b.itemToMove.shortName)
	default:
		inputView = ""
	}

	return bubbleStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			b.list.View(),
			inputStyle.Render(inputView),
		))
}
