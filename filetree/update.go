package filetree

import (
	"errors"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/knipferrc/teacup/dirfs"
	"github.com/muesli/termenv"
)

const (
	yesKey   = "y"
	enterKey = "enter"
)

// Update handles updating the filetree.
func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
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
		if b.IsFiltering() {
			break
		}

        if !b.active {
            return b, nil
        }

		switch b.state {
		case deleteItemState:
			if msg.String() == yesKey {
				selectedItem := b.GetSelectedItem()

				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully deleted item"),
				)

				cmds = append(cmds, statusCmd, tea.Sequentially(
					deleteItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden, b.showIcons),
				))

				b.state = idleState

				return b, tea.Batch(cmds...)
			}
		case moveItemState:
			if msg.String() == enterKey {
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully moved item"),
				)

				cmds = append(cmds, statusCmd, tea.Sequentially(
					moveItemCmd(b.itemToMove.path, b.itemToMove.shortName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden, b.showIcons),
				))

				b.state = idleState

				return b, tea.Batch(cmds...)
			}
		}

		switch {
		case key.Matches(msg, openDirectoryKey):
			if !b.input.Focused() {
				selectedDir := b.GetSelectedItem()
				cmds = append(cmds, getDirectoryListingCmd(selectedDir.fileName, b.showHidden, b.showIcons))
			}
		case key.Matches(msg, copyItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully copied file"),
				)

				cmds = append(cmds, statusCmd, tea.Sequentially(
					copyItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden, b.showIcons),
				))
			}
		case key.Matches(msg, zipItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully zipped item"),
				)

				cmds = append(cmds, statusCmd, tea.Sequentially(
					zipItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden, b.showIcons),
				))
			}
		case key.Matches(msg, unzipItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully unzipped item"),
				)

				cmds = append(cmds, statusCmd, tea.Sequentially(
					unzipItemCmd(selectedItem.fileName),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden, b.showIcons),
				))
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
				b.state = deleteItemState

				return b, nil
			}
		case key.Matches(msg, moveItemKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				b.state = moveItemState
				b.itemToMove = itemToMove{
					shortName: selectedItem.shortName,
					path:      selectedItem.fileName,
				}

				return b, nil
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
				cmds = append(cmds, getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden, b.showIcons))
			}
		case key.Matches(msg, homeShortcutKey):
			if !b.input.Focused() {
				cmds = append(cmds, getDirectoryListingCmd(dirfs.HomeDirectory, b.showHidden, b.showIcons))
			}
		case key.Matches(msg, rootShortcutKey):
			if !b.input.Focused() {
				cmds = append(cmds, getDirectoryListingCmd(dirfs.RootDirectory, b.showHidden, b.showIcons))
			}
		case key.Matches(msg, copyToClipboardKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()
				cmds = append(cmds, copyToClipboardCmd(selectedItem.fileName))
			}
		case key.Matches(msg, escapeKey):
			b.state = idleState

			if b.input.Focused() {
				b.input.Reset()
				b.input.Blur()
			}
		case key.Matches(msg, openInEditorKey):
			if !b.input.Focused() {
				selectedItem := b.GetSelectedItem()

				if b.selectionPath == "" && !selectedItem.IsDirectory() {
					editorPath := os.Getenv("EDITOR")
					if editorPath == "" {
						return b, handleErrorCmd(errors.New("$EDITOR not set"))
					}

					editorCmd := exec.Command(editorPath, selectedItem.FileName())
					editorCmd.Stdin = os.Stdin
					editorCmd.Stdout = os.Stdout
					editorCmd.Stderr = os.Stderr

					err := editorCmd.Run()
					termenv.AltScreen()

					if err != nil {
						return b, handleErrorCmd(err)
					}

					return b, tea.Batch(b.redrawCmd(), tea.HideCursor)
				}

				return b, tea.Sequentially(
					writeSelectionPathCmd(b.selectionPath, selectedItem.ShortName()),
					tea.Quit,
				)
			}
		case key.Matches(msg, submitInputKey):
			selectedItem := b.GetSelectedItem()

			switch b.state {
			case idleState, deleteItemState, moveItemState:
				return b, nil
			case createFileState:
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully created file"),
				)

				cmds = append(cmds, statusCmd, tea.Sequentially(
					createFileCmd(b.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden, b.showIcons),
				))
			case createDirectoryState:
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully created directory"),
				)

				cmds = append(cmds, statusCmd, tea.Sequentially(
					createDirectoryCmd(b.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden, b.showIcons),
				))
			case renameItemState:
				statusCmd := b.list.NewStatusMessage(
					statusMessageInfoStyle("Successfully renamed"),
				)

				cmds = append(cmds, statusCmd, tea.Sequentially(
					renameItemCmd(selectedItem.fileName, b.input.Value()),
					getDirectoryListingCmd(dirfs.CurrentDirectory, b.showHidden, b.showIcons),
				))
			}

			b.state = idleState
			b.input.Blur()
			b.input.Reset()
		}
	}

	if b.active {
		switch b.state {
		case idleState, moveItemState:
			b.list, cmd = b.list.Update(msg)
			cmds = append(cmds, cmd)
		case createFileState, createDirectoryState, renameItemState:
			b.input, cmd = b.input.Update(msg)
			cmds = append(cmds, cmd)
		case deleteItemState:
			return b, nil
		}
	}

	return b, tea.Batch(cmds...)
}
