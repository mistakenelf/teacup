// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import "github.com/charmbracelet/lipgloss"

// View returns a string representation of a filetree.
func (b Bubble) View() string {
	return listStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			inputStyle.Render(b.input.View()),
			b.list.View(),
		))
}
