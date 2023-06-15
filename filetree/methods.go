package filetree

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/dirfs"
)

const (
	thousand    = 1000
	ten         = 10
	fivePercent = 0.0499
)

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

// SetSize sets the size of the filetree.
func (m *Model) SetSize(width, height int) {
	horizontal, vertical := bubbleStyle.GetFrameSize()

	m.list.Styles.StatusBar.Width(width - horizontal)
	m.list.SetSize(
		width-horizontal-vertical,
		height-vertical-lipgloss.Height(m.input.View())-inputStyle.GetVerticalPadding(),
	)
}

// SetBorderColor sets the color of the border.
func (m *Model) SetBorderColor(color lipgloss.AdaptiveColor) {
	bubbleStyle = bubbleStyle.Copy().BorderForeground(color)
}

// GetSelectedItem returns the currently selected item in the tree.
func (m Model) GetSelectedItem() Item {
	selectedDir, ok := m.list.SelectedItem().(Item)
	if ok {
		return selectedDir
	}

	return Item{}
}

// Cursor returns the current position of the cursor in the tree.
func (m Model) Cursor() int {
	return m.list.Index() + 1
}

// TotalItems returns the total number of items in the tree.
func (m Model) TotalItems() int {
	return len(m.list.Items())
}

// SetIsActive sets if the bubble is currently active.
func (m *Model) SetIsActive(active bool) {
	m.active = active
}

// IsFiltering returns if the tree is currently being filtered.
func (m Model) IsFiltering() bool {
	return m.list.FilterState() == list.Filtering
}

// SetStartDir sets a starting directory.
func (m *Model) SetStartDir(dir string) {
	m.startDir = dir
}

// SetSelectionPath sets the path in which to write to a file when editing.
func (m *Model) SetSelectionPath(path string) {
	m.selectionPath = path
}

// SetTitleColors sets the background and foreground of the title.
func (m *Model) SetTitleColors(foreground, background lipgloss.AdaptiveColor) {
	m.list.Styles.Title = m.list.Styles.Title.Copy().
		Bold(true).
		Italic(true).
		Background(background).
		Foreground(foreground)
}

// SetSelectedItemColors sets the foreground of the selected item.
func (m *Model) SetSelectedItemColors(foreground lipgloss.AdaptiveColor) {
	m.delegate.Styles.SelectedTitle = m.delegate.Styles.SelectedTitle.Copy().
		Foreground(foreground).
		BorderLeftForeground(foreground)
	m.delegate.Styles.SelectedDesc = m.delegate.Styles.SelectedTitle.Copy()

	m.list.SetDelegate(m.delegate)
}

// SetBorderless sets weather or not to show the border.
func (m *Model) SetBorderless(borderless bool) {
	if borderless {
		bubbleStyle = bubbleStyle.Copy().BorderStyle(lipgloss.HiddenBorder())
	} else {
		bubbleStyle = bubbleStyle.Copy().BorderStyle(lipgloss.NormalBorder())
	}
}

// ToggleShowIcons sets weather or not to show icons.
func (m *Model) ToggleShowIcons(showIcons bool) tea.Cmd {
	m.showIcons = showIcons

	return getDirectoryListingCmd(dirfs.CurrentDirectory, m.showHidden, m.showIcons)
}

// ToggleHelp sets weather or not to show the help section.
func (m *Model) ToggleHelp(showHelp bool) {
	m.list.SetShowHelp(showHelp)
}
