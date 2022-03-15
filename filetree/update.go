// Package filetree implements a filetree bubble which can be used
// to navigate the filesystem and perform actions on files and directories.
package filetree

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/knipferrc/teacup/dirfs"
)

// Update handles updating the filetree.
func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case getDirectoryListingMsg:
		if msg != nil {
			cmd = b.list.SetItems(msg)
			cmds = append(cmds, cmd)
		}
	case copyToClipboardMsg:
		return b, b.list.NewStatusMessage(statusMessageInfoStyle(string(msg)))
	case errorMsg:
		return b, b.list.NewStatusMessage(statusMessageErrorStyle(msg.Error()))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, openDirectoryKey):
			if !b.input.Focused() {
				selectedDir := b.GetSelectedItem()
				cmds = append(cmds, getDirectoryListingCmd(selectedDir.fileName, b.showHidden))
			}
		case key.Matches(msg, copyItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully copied file"),
				)

				cmds = append(cmds, tea.Sequentially(
					copyItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))
				cmds = append(cmds, statusCmd)
			}
		case key.Matches(msg, zipItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully zipped item"),
				)

				cmds = append(cmds, tea.Sequentially(
					zipItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))
				cmds = append(cmds, statusCmd)
			}
		case key.Matches(msg, unzipItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully unzipped item"),
				)

				cmds = append(cmds, tea.Sequentially(
					unzipItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))
				cmds = append(cmds, statusCmd)
			}
		case key.Matches(msg, createFileKey):
			if !b.input.Focused() {
				b.input.Focus()
				b.input.Placeholder = "Enter name of new file"
				b.state = createFileState

				return b, textinput.Blink
			}
		case key.Matches(msg, createDirectoryKey):
			if !b.input.Focused() {
				b.input.Focus()
				b.input.Placeholder = "Enter name of new directory"
				b.state = createDirectoryState

				return b, textinput.Blink
			}
		case key.Matches(msg, deleteItemKey):
			if !b.input.Focused() {
				b.input.Focus()
				b.input.Placeholder = "Are you sure you want to delete (y/n)?"
				b.state = deleteItemState

				return b, textinput.Blink
			}
		case key.Matches(msg, renameItemKey):
			if !b.input.Focused() {
				b.input.Focus()
				b.input.Placeholder = "Enter new name"
				b.state = renameItemState

				return b, textinput.Blink
			}
		case key.Matches(msg, toggleHiddenKey):
			if !b.input.Focused() {
				b.showHidden = !b.showHidden
				cmds = append(cmds, getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden))
			}
		case key.Matches(msg, homeShortcutKey):
			if !b.input.Focused() {
				cmds = append(cmds, getDirectoryListingCmd(dirfs.HomeDirectory, b.showHidden))
			}
		case key.Matches(msg, copyToClipboardKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				cmds = append(cmds, copyToClipboardCmd(selectedItem.fileName))
			}
		case key.Matches(msg, escapeKey):
			if b.input.Focused() {
				b.input.Reset()
				b.input.Blur()
				b.state = idleState
			}
		case key.Matches(msg, submitInputKey):
			switch b.state {
			case idleState:
				return b, nil
			case createFileState:
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully created file"),
				)

				cmds = append(cmds, tea.Sequentially(
					createFileCmd(b.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))
				cmds = append(cmds, statusCmd)

				b.state = idleState
				b.input.Blur()
				b.input.Reset()
			case createDirectoryState:
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully created directory"),
				)

				cmds = append(cmds, statusCmd)
				cmds = append(cmds, tea.Sequentially(
					createDirectoryCmd(b.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))

				b.state = idleState
				b.input.Blur()
				b.input.Reset()
			case deleteItemState:
				if strings.ToLower(b.input.Value()) == "y" {
					selectedDir := b.GetSelectedItem()

					statusCmd := b.list.NewStatusMessage(
						statusMessageInfoStyle("Successfully deleted item"),
					)

					cmds = append(cmds, statusCmd)
					cmds = append(cmds, tea.Sequentially(
						deleteItemCmd(selectedDir.fileName),
						getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
					))
				}

				b.state = idleState
				b.input.Blur()
				b.input.Reset()
			case renameItemState:
				statusCmd := b.list.NewStatusMessage(statusMessageInfoStyle("Successfully renamed"))
				cmds = append(cmds, statusCmd)

				selectedItem := b.GetSelectedItem()
				cmds = append(cmds, tea.Sequentially(
					renameItemCmd(selectedItem.fileName, b.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden),
				))

				b.state = idleState
				b.input.Blur()
				b.input.Reset()
			}
		}
	}

	switch b.state {
	case idleState:
		b.list, cmd = b.list.Update(msg)
		cmds = append(cmds, cmd)
	case createFileState, createDirectoryState, deleteItemState, renameItemState:
		b.input, cmd = b.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return b, tea.Batch(cmds...)
}
