// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/icons"
)

// item represents a list item.
type item struct {
	title            string
	desc             string
	fileName         string
	shortName        string
	extension        string
	currentDirectory string
	isDirectory      bool
	showIcons        bool
	fileInfo         fs.FileInfo
}

// Title returns the title of the list item.
func (i item) Title() string {
	if i.fileInfo != nil {
		icon, color := icons.GetIcon(
			i.fileInfo.Name(),
			filepath.Ext(i.fileInfo.Name()),
			icons.GetIndicator(i.fileInfo.Mode()),
		)
		fileIcon := lipgloss.NewStyle().Width(fileIconWidth).Render(fmt.Sprintf("%s%s\033[0m ", color, icon))

		if i.showIcons {
			return fmt.Sprintf("%s %s", i.title, fileIcon)
		}

		return i.title
	}

	return i.title
}

// FileName returns the file name of the list item.
func (i item) FileName() string { return i.fileName }

// FileExtension returns the extension of the list item.
func (i item) FileExtension() string { return i.extension }

// IsDirectory returns true if the list item is a directory.
func (i item) IsDirectory() bool { return i.isDirectory }

// Description returns the description of the list item.
func (i item) Description() string { return i.desc }

// FilterValue returns the current filter value.
func (i item) FilterValue() string { return i.title }

// ShortName returns the short name of the selected item.
func (i item) ShortName() string { return i.shortName }

// CurrentDirectory returns the current directory of the tree
func (i item) CurrentDirectory() string { return i.currentDirectory }
