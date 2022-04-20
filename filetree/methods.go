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
func (b *Bubble) SetSize(width, height int) {
	horizontal, vertical := bubbleStyle.GetFrameSize()

	b.list.Styles.StatusBar.Width(width - horizontal)
	b.list.SetSize(
		width-horizontal-vertical,
		height-vertical-lipgloss.Height(b.input.View())-inputStyle.GetVerticalPadding(),
	)
}

// SetBorderColor sets the color of the border.
func (b *Bubble) SetBorderColor(color lipgloss.AdaptiveColor) {
	bubbleStyle = bubbleStyle.Copy().BorderForeground(color)
}

// GetSelectedItem returns the currently selected item in the tree.
func (b Bubble) GetSelectedItem() Item {
	selectedDir, ok := b.list.SelectedItem().(Item)
	if ok {
		return selectedDir
	}

	return Item{}
}

// Cursor returns the current position of the cursor in the tree.
func (b Bubble) Cursor() int {
	return b.list.Index() + 1
}

// TotalItems returns the total number of items in the tree.
func (b Bubble) TotalItems() int {
	return len(b.list.Items())
}

// SetIsActive sets if the bubble is currently active.
func (b *Bubble) SetIsActive(active bool) {
	b.active = active
}

// IsFiltering returns if the tree is currently being filtered.
func (b Bubble) IsFiltering() bool {
	return b.list.FilterState() == list.Filtering
}

// SetStartDir sets a starting directory.
func (b *Bubble) SetStartDir(dir string) {
	b.startDir = dir
}

// SetSelectionPath sets the path in which to write to a file when editing.
func (b *Bubble) SetSelectionPath(path string) {
	b.selectionPath = path
}

// SetTitleColors sets the background and foreground of the title.
func (b *Bubble) SetTitleColors(foreground, background lipgloss.AdaptiveColor) {
	b.list.Styles.Title = b.list.Styles.Title.Copy().
		Bold(true).
		Italic(true).
		Background(background).
		Foreground(foreground)
}

// SetSelectedItemColors sets the foreground of the selected item.
func (b *Bubble) SetSelectedItemColors(foreground lipgloss.AdaptiveColor) {
	b.delegate.Styles.SelectedTitle = b.delegate.Styles.SelectedTitle.Copy().
		Foreground(foreground).
		BorderLeftForeground(foreground)
	b.delegate.Styles.SelectedDesc = b.delegate.Styles.SelectedTitle.Copy()

	b.list.SetDelegate(b.delegate)
}

// SetBorderless sets weather or not to show the border.
func (b *Bubble) SetBorderless(borderless bool) {
	if borderless {
		bubbleStyle = bubbleStyle.Copy().BorderStyle(lipgloss.HiddenBorder())
	} else {
		bubbleStyle = bubbleStyle.Copy().BorderStyle(lipgloss.NormalBorder())
	}
}

// ToggleShowIcons sets weather or not to show icons.
func (b *Bubble) ToggleShowIcons(showIcons bool) tea.Cmd {
	b.showIcons = showIcons

	return getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden, b.showIcons)
}

// ToggleHelp sets weather or not to show the help section.
func (b *Bubble) ToggleHelp(showHelp bool) {
	b.list.SetShowHelp(showHelp)
}
