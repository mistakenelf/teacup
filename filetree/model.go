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
	renameItemState
	moveItemState
)

type itemToMove struct {
	shortName string
	path      string
}

// Bubble represents the properties of a filetree.
type Model struct {
	state         sessionState
	list          list.Model
	input         textinput.Model
	showHidden    bool
	showIcons     bool
	active        bool
	width         int
	height        int
	startDir      string
	selectionPath string
	itemToMove    itemToMove
	delegate      list.DefaultDelegate
}

// New creates a new instance of a filetree.
func New(
	active, borderless bool,
	startDir, selectionPath string,
	borderColor, selectedItemColor, titleBackgroundColor, titleForegroundColor lipgloss.AdaptiveColor,
) Model {
	listDelegate := list.NewDefaultDelegate()
	listDelegate.Styles.SelectedTitle = listDelegate.Styles.SelectedTitle.Copy().
		Foreground(selectedItemColor).
		BorderLeftForeground(selectedItemColor)
	listDelegate.Styles.SelectedDesc = listDelegate.Styles.SelectedTitle.Copy()

	listModel := list.New([]list.Item{}, listDelegate, 0, 0)
	listModel.Title = "Filetree"
	listModel.Styles.Title = listModel.Styles.Title.Copy().
		Bold(true).
		Italic(true).
		Background(titleBackgroundColor).
		Foreground(titleForegroundColor)
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
			renameItemKey,
			openInEditorKey,
			submitInputKey,
			moveItemKey,
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
			renameItemKey,
			openInEditorKey,
			submitInputKey,
			moveItemKey,
		}
	}

	input := textinput.New()
	input.Prompt = "❯ "
	input.Placeholder = "Enter file name"
	input.CharLimit = 250
	input.Width = 50

	if borderless {
		bubbleStyle = bubbleStyle.Copy().Border(lipgloss.HiddenBorder())
	} else {
		bubbleStyle = bubbleStyle.Copy().BorderForeground(borderColor)
	}

	return Model{
		list:          listModel,
		input:         input,
		showHidden:    true,
		showIcons:     true,
		active:        active,
		state:         idleState,
		startDir:      startDir,
		selectionPath: selectionPath,
		delegate:      listDelegate,
	}
}
