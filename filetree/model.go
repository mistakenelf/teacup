// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	idleState sessionState = iota
	createFileState
	createDirectoryState
	deleteItemState
)

// Bubble represents the properties of a filetree.
type Bubble struct {
	state      sessionState
	list       list.Model
	input      textinput.Model
	showHidden bool
}

// Item represents a list item.
type Item struct {
	ItemTitle string
	Desc      string
	FileName  string
}

// Title returns the title of the list item.
func (i Item) Title() string {
	return i.ItemTitle
}

// Description returns the description of the list item.
func (i Item) Description() string { return i.Desc }

// FilterValue returns the current filter value.
func (i Item) FilterValue() string { return i.ItemTitle }

// New creates a new instance of a filetree.
func New(borderColor lipgloss.AdaptiveColor, borderless bool) Bubble {
	listModel := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	listModel.Title = "Filetree"
	listModel.DisableQuitKeybindings()
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
			homeShortcutKey,
			copyToClipboardKey,
			escapeKey,
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
			homeShortcutKey,
			copyToClipboardKey,
			escapeKey,
			submitInputKey,
		}
	}

	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.Placeholder = "Enter file name"
	input.CharLimit = 250
	input.Width = 50

	if borderless {
		bubbleStyle = bubbleStyle.Copy().Border(lipgloss.HiddenBorder())
	} else {
		bubbleStyle = bubbleStyle.Copy().BorderForeground(borderColor)
	}

	return Bubble{
		list:       listModel,
		input:      input,
		showHidden: true,
		state:      idleState,
	}
}
