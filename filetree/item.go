// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

// item represents a list item.
type item struct {
	title       string
	desc        string
	fileName    string
	extension   string
	isDirectory bool
}

// Title returns the title of the list item.
func (i item) Title() string {
	return i.title
}

// GetFileName returns the file name of the list item.
func (i item) GetFileName() string { return i.fileName }

// GetFileExtension returns the extension of the list item.
func (i item) GetFileExtension() string { return i.extension }

// IsDirectory returns true if the list item is a directory.
func (i item) IsDirectory() bool { return i.isDirectory }

// Description returns the description of the list item.
func (i item) Description() string { return i.desc }

// FilterValue returns the current filter value.
func (i item) FilterValue() string { return i.title }
