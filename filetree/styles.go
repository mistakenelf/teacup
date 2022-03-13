// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import "github.com/charmbracelet/lipgloss"

const (
	fileIconWidth = 2
)

var (
	bubbleStyle = lipgloss.NewStyle().
			Padding(1).
			BorderStyle(lipgloss.NormalBorder())
	inputStyle             = lipgloss.NewStyle().PaddingBottom(1)
	statusMessageInfoStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	statusMessageErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}).
				Render
)
