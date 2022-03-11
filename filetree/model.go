// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type appState int

const (
	createFileState appState = iota
	createDirectoryState
	deleteItemState
)

// Bubble represents the properties of a filetree.
type Bubble struct {
	state      appState
	list       list.Model
	input      textinput.Model
	showHidden bool
}

// item represents a list item.
type item struct {
	title    string
	desc     string
	fileName string
}

// Title returns the title of the list item.
func (i item) Title() string {
	return i.title
}

// Description returns the description of the list item.
func (i item) Description() string { return i.desc }

// FilterValue returns the current filter value.
func (i item) FilterValue() string { return i.title }

// New creates a new instance of a filetree.
func New(borderColor lipgloss.AdaptiveColor, borderless bool) Bubble {
	listModel := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	listModel.Title = "Filetree"
	listModel.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			openDirectoryKey,
			createFileKey,
			createDirectoryKey,
			deleteItemKey,
			copyItemKey,
			zipItemKey,
			unzipItemKey,
			toggleHiddenKey,
			submitInputKey,
		}
	}
	listModel.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			openDirectoryKey,
			createFileKey,
			createDirectoryKey,
			deleteItemKey,
			copyItemKey,
			zipItemKey,
			unzipItemKey,
			toggleHiddenKey,
			submitInputKey,
		}
	}

	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.Placeholder = "Enter file name"
	input.CharLimit = 250
	input.Width = 50

	if borderless {
		listStyle = listStyle.Copy().Border(lipgloss.HiddenBorder())
	} else {
		listStyle = listStyle.Copy().BorderForeground(borderColor)
	}

	return Bubble{
		list:       listModel,
		input:      input,
		showHidden: true,
	}
}
