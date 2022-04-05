// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import "github.com/charmbracelet/bubbles/key"

var (
	openDirectoryKey   = key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "open directory"))
	createFileKey      = key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "create file"))
	submitInputKey     = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit input value"))
	createDirectoryKey = key.NewBinding(key.WithKeys("N"), key.WithHelp("N", "create directory"))
	deleteItemKey      = key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "delete item"))
	copyItemKey        = key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy item"))
	zipItemKey         = key.NewBinding(key.WithKeys("z"), key.WithHelp("z", "zip item"))
	unzipItemKey       = key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "unzip item"))
	toggleHiddenKey    = key.NewBinding(key.WithKeys("."), key.WithHelp(".", "toggle hidden files"))
	homeShortcutKey    = key.NewBinding(key.WithKeys("~"), key.WithHelp("~", "go to home directory"))
	rootShortcutKey    = key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "go to root directory"))
	copyToClipboardKey = key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "copy path to clipboard"))
	renameItemKey      = key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rename item"))
	openInEditorKey    = key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "open in editor"))
	escapeKey          = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "reset to initial state"))
)
