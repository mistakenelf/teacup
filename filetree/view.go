// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import "github.com/charmbracelet/lipgloss"

// View returns a string representation of a filetree.
func (b Bubble) View() string {
	var inputView string

	switch b.state {
	case idleState:
		inputView = ""
	case createFileState, createDirectoryState, deleteItemState, renameItemState:
		inputView = b.input.View()
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
